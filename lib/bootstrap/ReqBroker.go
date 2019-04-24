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
	GetSnapshot() bool
}

type ReqBrokerOptions interface {
	GetVersion() string
}

type ReqBroker struct {
	httpInvoker *engine.HttpInvoker
}

type SnapshotOutput struct {}

func (g *SnapshotOutput) GetTargetWriter() io.Writer {
	return os.Stdout
}

type ExplanationWriter struct {}

func (z *ExplanationWriter) GetConsoleOut() io.Writer {
	return os.Stdout
}

func (z *ExplanationWriter) GetConsoleErr() io.Writer {
	return os.Stderr
}

func NewReqBroker(opts ReqBrokerOptions) (obj *ReqBroker, err error) {
	obj = &ReqBroker{}
	httpInvokerOptions := &engine.HttpInvokerOptions{}
	if opts != nil {
		httpInvokerOptions.Version = opts.GetVersion()
	}
	obj.httpInvoker, err = engine.NewHttpInvoker(httpInvokerOptions)
	if err != nil {
		return nil, err
	}
	return obj, err
}

func (z *ReqBroker) Execute(args ReqArguments) error {
	z.assertReady()
	if args == nil {
		return fmt.Errorf("ReqBroker.Execute() arguments must not be nil")
	}

	if args.GetSnapshot() {
		_, err := z.httpInvoker.Do(transformReqArgs(args), &SnapshotOutput{})
		if err != nil {
			return err
		}
		return nil
	}

	_, err := z.httpInvoker.Do(transformReqArgs(args), &ExplanationWriter{})
	if err != nil {
		return err
	}
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
