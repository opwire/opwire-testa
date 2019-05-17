package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/opwire/opwire-testa/lib/utils"
)

type HttpInvoker interface {
	Do(req *HttpRequest, interceptors ...Interceptor) (res *HttpResponse, err error)
}

type HttpInvokerOptions struct {
	PDP string
}

type HttpInvokerImpl struct {
	pdp string
}

func NewHttpInvoker(opts *HttpInvokerOptions) (c *HttpInvokerImpl, err error) {
	c = &HttpInvokerImpl{}
	if opts != nil {
		c.pdp = opts.PDP
	}
	return c, nil
}

func (c *HttpInvokerImpl) Do(req *HttpRequest, interceptors ...Interceptor) (*HttpResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("Request must not be nil")
	}

	if len(req.PDP) == 0 {
		req.PDP = c.pdp
	}

	var reqTimeout time.Duration
	if req.Timeout != nil {
		var err error
		reqTimeout, err = time.ParseDuration(*req.Timeout)
		if err != nil || reqTimeout <= 0 {
			reqTimeout = time.Second * 10
		}
	}

	var httpClient *http.Client = &http.Client{
		Timeout: reqTimeout,
	}

	lowReq, err := req.GetRawRequest()
	if err != nil {
		return nil, err
	}

	// Pre-processing
	for _, interceptor := range interceptors {
		if processor, ok := interceptor.(PreProcessor); processor != nil && ok {
			processor.PreProcess(req)
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
		if processor, ok := interceptor.(PostProcessor); processor != nil && ok {
			processor.PostProcess(req, res)
		}
	}

	return res, nil
}

func BuildUrl(req *HttpRequest) string {
	url := req.Url
	if len(url) == 0 {
		pdp := utils.DEFAULT_PDP
		if len(req.PDP) > 0 {
			pdp = req.PDP
		}
		basePath := utils.DEFAULT_PATH
		if len(req.Path) > 0 {
			basePath = req.Path
		}
		url, _ = utils.UrlJoin(pdp, basePath)
	}
	return url
}

type HttpHeader struct {
	Name string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type HttpRequest struct {
	Method string `yaml:"method,omitempty" json:"method"`
	Url string `yaml:"url,omitempty" json:"url"`
	PDP string `yaml:"pdp,omitempty" json:"pdp"`
	Path string `yaml:"path,omitempty" json:"path"`
	Headers []HttpHeader `yaml:"headers,omitempty" json:"headers"`
	Body string `yaml:"body,omitempty" json:"body"`
	Timeout *string `yaml:"timeout,omitempty" json:"timeout"`
	request *http.Request
}

func (r *HttpRequest) GetRawRequest() (req *http.Request, err error) {
	if r.request == nil {
		url := BuildUrl(r)

		method := "GET"
		if len(r.Method) > 0 {
			method = r.Method
		}

		var body *bytes.Buffer

		if len(r.Body) > 0 {
			body = bytes.NewBufferString(r.Body)
		} else {
			body = bytes.NewBuffer([]byte{})
		}

		req, err = http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}

		for _, header := range r.Headers {
			if len(header.Name) > 0 && len(header.Value) > 0 {
				req.Header.Add(header.Name, header.Value)
			}
		}

		r.request = req
	}
	return r.request, nil
}

type HttpResponse struct {
	Status string
	StatusCode int
	Version string
	Header http.Header
	ContentLength int64
	Body []byte
}

type Interceptor interface {}

type PreProcessor interface {
	Interceptor
	PreProcess(req *HttpRequest) error
}

type PostProcessor interface {
	Interceptor
	PostProcess(req *HttpRequest, res *HttpResponse) error
}
