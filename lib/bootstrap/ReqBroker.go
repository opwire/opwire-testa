package bootstrap

import (
	"fmt"
	"io"
	"os"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/utils"
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

	res, err := z.httpInvoker.Do(transformReqArgs(args), z)
	if err != nil {
		return err
	}

	_ = res

	return nil
}

func (z *ReqBroker) GetConsoleOut() io.Writer {
	return z.consoleOut
}

func (z *ReqBroker) GetConsoleErr() io.Writer {
	return z.consoleErr
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
			pair := utils.Split(item, ":")
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
