package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/opwire/opwire-testa/lib/utils"
)

type HttpInvokerOptions interface {
	GetPDP() string
}

type HttpInvoker struct {
	pdp string
}

func NewHttpInvoker(opts HttpInvokerOptions) (*HttpInvoker, error) {
	c := &HttpInvoker{}
	if opts != nil {
		c.pdp = opts.GetPDP()
	}
	if len(c.pdp) == 0 {
		c.pdp = DEFAULT_PDP
	}
	return c, nil
}

func (c *HttpInvoker) Do(req *HttpRequest, interceptors ...Interceptor) (*HttpResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("Request must not be nil")
	}

	url := req.Url
	if len(url) == 0 {
		pdp := c.pdp
		if len(req.PDP) > 0 {
			pdp = req.PDP
		}
		basePath := "/$"
		if len(req.Path) > 0 {
			basePath = req.Path
		}
		url, _ = utils.UrlJoin(pdp, basePath)
	}

	reqTimeout := time.Second * 10
	var httpClient *http.Client = &http.Client{
		Timeout: reqTimeout,
	}

	method := "GET"
	if len(req.Method) > 0 {
		method = req.Method
	}

	var body *bytes.Buffer
	
	if len(req.Body) > 0 {
		body = bytes.NewBufferString(req.Body)
	}
	
	lowReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for _, header := range req.Headers {
		if len(header.Name) > 0 && len(header.Value) > 0 {
			lowReq.Header.Add(header.Name, header.Value)
		}
	}

	// Pre-processing
	for _, interceptor := range interceptors {
		if monitor, ok := interceptor.(ExplanationWriter); monitor != nil && ok {
			w := monitor.GetConsoleOut()
			if w != nil {
				renderRequest(w, lowReq)
			}
		}
	}

	// Make HTTP request
	lowRes, err := httpClient.Do(lowReq)
	if lowRes != nil && lowRes.Body != nil {
		defer lowRes.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	res := &HttpResponse{}

	res.Version = lowRes.Proto
	res.Status = lowRes.Status
	res.StatusCode = lowRes.StatusCode
	res.Header = lowRes.Header

	res.Body, err = ioutil.ReadAll(lowRes.Body)
	if err != nil {
		return nil, err
	}

	// Post-processing
	for _, interceptor := range interceptors {
		if monitor, ok := interceptor.(ExplanationWriter); monitor != nil && ok {
			w := monitor.GetConsoleOut()
			if w != nil {
				renderResponse(w, res)
			}
		}
		if snapshot, ok := interceptor.(SnapshotGenerator); snapshot != nil && ok {
			w := snapshot.GetTargetWriter()
			if w != nil {
				generateTestCase(w, req, res)
			}
		}
	}

	return res, nil
}

func renderRequest(w io.Writer, req *http.Request) error {
	// render first line
	line := []string{">"}
	if len(req.Method) > 0 {
		line = append(line, req.Method)
	}
	if req.URL != nil {
		reqURI := req.URL.RequestURI()
		if len(reqURI) > 0 {
			line = append(line, reqURI)
		} else {
			if len(req.URL.Path) > 0 {
				line = append(line, req.URL.Path)
			}
		}
	}
	if len(req.Proto) > 0 {
		line = append(line, req.Proto)
	}
	fmt.Fprintln(w, strings.Join(line, " "))
	// render Host
	if req.URL != nil && len(req.URL.Host) > 0 {
		fmt.Fprintln(w, "> Host: " + req.URL.Host)
	}
	// render User-Agent
	userAgent := req.UserAgent()
	if len(userAgent) > 0 {
		fmt.Fprintln(w, "> User-Agent: " + userAgent)
	}
	// render headers
	for key, vals := range req.Header {
		for _, val := range vals {
			fmt.Fprintln(w, "> " + key + ": " + val)
		}
	}
	fmt.Fprintln(w, ">")
	return nil
}

func renderResponse(w io.Writer, res *HttpResponse) error {
	// render status line
	line := []string{"<"}
	if len(res.Version) > 0 {
		line = append(line, res.Version)
	}
	if len(res.Status) > 0 {
		line = append(line, res.Status)
	} else {
		line = append(line, fmt.Sprintf("%v", res.StatusCode))
	}
	fmt.Fprintln(w, strings.Join(line, " "))
	// render headers
	for key, vals := range res.Header {
		for _, val := range vals {
			fmt.Fprintln(w, "< " + key + ": " + val)
		}
	}
	fmt.Fprintln(w, "<")
	// render body
	fmt.Fprintln(w, string(res.Body))
	return nil
}

func generateTestCase(w io.Writer, req *HttpRequest, res *HttpResponse) error {
	s := TestCase{}
	s.Title = "<Generated testcase>"
	s.Request = req
	s.Expectation = generateExpectation(res)

	r := &GeneratedSnapshot{}
	r.TestCases = []TestCase{s}
	script, err := yaml.Marshal(r)
	if err != nil {
		fmt.Fprintln(w, "Cannot marshal generated testcase, error: %s", err)
		return err
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, string(script))

	return nil
}

func generateExpectation(res *HttpResponse) *Expectation {
	if res == nil {
		return nil
	}
	e := &Expectation{}

	// status-code
	sc := res.StatusCode
	e.StatusCode = &MeasureStatusCode{
		IsEqualTo: &sc,
	}

	// header
	total := len(res.Header)
	if total > 0 {
		e.Headers = &MeasureHeaders{
			HasTotal: &total,
			Items: make([]MeasureHeader, 0),
		}
		for key, vals := range res.Header {
			if len(vals) == 1 {
				name := key
				value := vals[0]
				one := MeasureHeader{
					Name: &name,
					IsEqualTo: &value,
				}
				e.Headers.Items = append(e.Headers.Items, one)
			}
		}
	}

	// body
	e.Body = &MeasureBody{}

	obj := make(map[string]interface{}, 0)
	if e.Body.HasFormat == nil {
		if err := json.Unmarshal(res.Body, &obj); err == nil {
			e.Body.HasFormat = RefOfString("json")
			var content string
			if out, err := json.MarshalIndent(obj, "", "  "); err == nil {
				content = string(out)
			} else {
				content = string(res.Body)
			}
			e.Body.Includes = &content
		}
	}

	if e.Body.HasFormat == nil {
		if err := yaml.Unmarshal(res.Body, &obj); err == nil {
			e.Body.HasFormat = RefOfString("yaml")
			var content string
			if out, err := yaml.Marshal(obj); err == nil {
				content = string(out)
			} else {
				content = string(res.Body)
			}
			e.Body.Includes = &content
		}
	}

	if e.Body.HasFormat == nil {
		e.Body.HasFormat = RefOfString("flat")
		e.Body.IsEqualTo = RefOfString(string(res.Body))
		e.Body.MatchWith = RefOfString(".*")
	}

	return e
}

func RefOfString(val string) *string {
	return &val
}

const DEFAULT_PDP string = `http://localhost:17779`

type HttpHeader struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}

type HttpRequest struct {
	Method string `yaml:"method,omitempty"`
	Url string `yaml:"url,omitempty"`
	PDP string `yaml:"pdp,omitempty"`
	Path string `yaml:"path,omitempty"`
	Headers []HttpHeader `yaml:"headers,omitempty"`
	Body string `yaml:"body,omitempty"`
}

type HttpResponse struct {
	Status string
	StatusCode int
	Version string
	Header http.Header
	ContentLength int64
	Body []byte
}

type Interceptor interface {
}

type ExplanationWriter interface {
	Interceptor
	GetConsoleOut() io.Writer
	GetConsoleErr() io.Writer
}

type SnapshotGenerator interface {
	Interceptor
	GetTargetWriter() io.Writer
}

type GeneratedSnapshot struct {
	TestCases []TestCase `yaml:"testcase-snapshot"`
}
