package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/opwire/opwire-testa/lib/engine"
)

type ReqArguments interface {
	GetMethod() string
	GetUrl() string
	GetHeader() []string
	GetBody() string
}

type ReqBrokerOptions interface {}

type ReqBroker struct {
	options ReqBrokerOptions
	httpInvoker *engine.HttpInvoker
	consoleOut io.Writer
	consoleErr io.Writer
}

func NewReqBroker(opts ReqBrokerOptions) (obj *ReqBroker, err error) {
	obj = &ReqBroker{}
	obj.httpInvoker, err = engine.NewHttpInvoker(nil)
	if err != nil {
		return nil, err
	}
	obj.consoleOut = os.Stdout
	obj.consoleErr = os.Stderr
	return obj, err
}

func (z *ReqBroker) Execute(args ReqArguments) error {
	z.assertReady()
	if args == nil {
		return fmt.Errorf("ReqBroker.Execute() arguments must not be nil")
	}
	
	res, err := z.httpInvoker.Do(transformReqArgs(args))
	if err != nil {
		return err
	}

	z.displayResponse(res)

	return nil
}

func (z *ReqBroker) displayResponse(res *engine.HttpResponse) error {
	line := []string{ "<" }
	if len(res.Version) > 0 {
		line = append(line, res.Version)
	}
	if len(res.Status) > 0 {
		line = append(line, res.Status)
	} else {
		line = append(line, fmt.Sprintf("%v", res.StatusCode))
	}
	fmt.Fprintln(z.consoleOut, strings.Join(line, " "))
	return nil
}

func transformReqArgs(args ReqArguments) *engine.HttpRequest {
	req := &engine.HttpRequest{}

	req.Method = args.GetMethod()
	if len(req.Method) == 0 {
		req.Method = "GET"
	}

	req.Url = args.GetUrl()
	if len(req.Url) == 0 {
		req.Url = "http://localhost:17779/$"
	}

	req.Headers = make([]engine.HttpHeader, 0)
	headerList := args.GetHeader()
	if headerList != nil {
		for _, item := range headerList {
			pair := strings.Split(item, "=")
			if len(pair) == 2 {
				header := engine.HttpHeader{
					Name: pair[0],
					Value: pair[1],
				}
				req.Headers = append(req.Headers, header)
			}
		}
	}

	req.Body = args.GetBody()

	return req
}

func (z *ReqBroker) assertReady() {
	if z.httpInvoker == nil {
		panic(fmt.Errorf("httpInvoker must not be nil"))
	}
}
