package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/format"
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
	GetNoColor() bool
}

type ReqController struct {
	httpInvoker *engine.HttpInvoker
	outputPrinter *format.OutputPrinter
}

type GenerationPrinter struct {}

func (g *GenerationPrinter) GetWriter() io.Writer {
	return os.Stdout
}

type ExplanationTarget struct {}

func (z *ExplanationTarget) GetConsoleOut() io.Writer {
	return os.Stdout
}

func (z *ExplanationTarget) GetConsoleErr() io.Writer {
	return os.Stderr
}

func NewReqController(opts ReqControllerOptions) (obj *ReqController, err error) {
	obj = &ReqController{}

	// create a HTTP Invoker instance
	httpInvokerOptions := &engine.HttpInvokerOptions{}
	if opts != nil {
		httpInvokerOptions.Version = opts.GetVersion()
	}
	obj.httpInvoker, err = engine.NewHttpInvoker(httpInvokerOptions)
	if err != nil {
		return nil, err
	}

	// create a OutputPrinter instance
	obj.outputPrinter, err = format.NewOutputPrinter(opts)
	if err != nil {
		return nil, err
	}

	return obj, err
}

func (z *ReqController) Execute(args ReqArguments) error {
	z.assertReady(args)

	console := &ExplanationTarget{}

	if args.GetFormat() == "testcase" {
		_, err := z.httpInvoker.Do(transformReqArgs(args), &GenerationPrinter{})
		if err != nil {
			return z.displayError(err, console)
		}
		return nil
	}

	res, err := z.httpInvoker.Do(transformReqArgs(args), console)
	if err != nil {
		return z.displayError(err, console)
	}
	z.displayResult(res, console)
	return nil
}

func (z *ReqController) displayError(err error, console *ExplanationTarget) error {
	if err == nil || console == nil {
		return err
	}
	w := console.GetConsoleErr()
	if w == nil {
		return err
	}
	fmt.Fprintf(w, "* %s\n", err.Error())
	return err
}

func (z *ReqController) displayResult(res *engine.HttpResponse, console *ExplanationTarget) *engine.HttpResponse {
	if res == nil || console == nil {
		return res
	}
	w := console.GetConsoleOut()
	if w == nil {
		return res
	}
	var codeStr string
	if res.StatusCode < 400 {
		codeStr = z.outputPrinter.InfoMsg("StatusCode [" + strconv.Itoa(res.StatusCode) + "]")
	} else {
		codeStr = z.outputPrinter.WarnMsg("StatusCode [" + strconv.Itoa(res.StatusCode) + "]")
	}
	fmt.Fprintf(w, "* Please make sure %s is your expected result\n", codeStr)
	return res
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
