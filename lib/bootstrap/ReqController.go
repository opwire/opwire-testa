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
	GetFormat() string
}

type ReqControllerOptions interface {
	GetVersion() string
}

type ReqController struct {
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

func NewReqController(opts ReqControllerOptions) (obj *ReqController, err error) {
	obj = &ReqController{}
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

func (z *ReqController) Execute(args ReqArguments) error {
	z.assertReady(args)

	if args.GetFormat() == "testcase" {
		_, err := z.httpInvoker.Do(transformReqArgs(args), &SnapshotOutput{})
		if err != nil {
			return displayError(err)
		}
		return nil
	}

	_, err := z.httpInvoker.Do(transformReqArgs(args), &ExplanationWriter{})
	if err != nil {
		return displayError(err)
	}
	return nil
}

func displayError(err error) error {
	fmt.Printf("* %s\n", err.Error())
	return err
}

func transformReqArgs(args ReqArguments) *engine.HttpRequest {
	req := &engine.HttpRequest{}

	req.Method = args.GetMethod()
	if len(req.Method) == 0 {
		req.Method = "GET"
	}

	req.Url = args.GetUrl()
	if len(req.Url) == 0 {
		req.Url, _ = utils.UrlJoin(utils.DEFAULT_PDP, utils.DEFAULT_PATH)
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

func (z *ReqController) assertReady(a ...interface{}) {
	if z.httpInvoker == nil {
		panic(fmt.Errorf("httpInvoker must not be nil"))
	}
	if len(a) > 0 {
		for i, s := range a {
			if s == nil {
				panic(fmt.Errorf("argument[%d] must not be nil", i))
			}
		}
	}
}
