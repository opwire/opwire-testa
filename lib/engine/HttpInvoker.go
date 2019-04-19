package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

func (c *HttpInvoker) Do(req *HttpRequest) (*HttpResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("Request must not be nil")
	}
	pdp := c.pdp
	if len(req.PDP) > 0 {
		pdp = req.PDP
	}
	basePath := "/$"
	if len(req.Path) > 0 {
		basePath = req.Path
	}
	url, _ := utils.UrlJoin(pdp, basePath)

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

	lowRes, err := httpClient.Do(lowReq)
	if lowRes != nil && lowRes.Body != nil {
		defer lowRes.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	res := &HttpResponse{}

	res.Headers = make([]HttpHeader, 0)

	for key, _ := range lowRes.Header {
		val := lowRes.Header.Get(key)
		if len(val) > 0 {
			res.Headers = append(res.Headers, HttpHeader{
				Name: key,
				Value: val,
			})
		}
	}

	res.Body, err = ioutil.ReadAll(lowRes.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

const DEFAULT_PDP string = `http://localhost:17779`

type HttpHeader struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}

type HttpRequest struct {
	Method string `yaml:"method"`
	PDP string `yaml:"pdp"`
	Path string `yaml:"path"`
	Headers []HttpHeader `yaml:"headers"`
	Body string `yaml:"body"`
}

type HttpResponse struct {
	Status string
	StatusCode int
	Version string
	Headers []HttpHeader
	ContentLength int64
	Body []byte
}

type HttpMeasure struct {
	Headers []HttpHeader `yaml:"headers"`
	Body string `yaml:"body"`
}
