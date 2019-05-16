package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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
	httpInvoker engine.HttpInvoker
	specGenerator *engine.SpecGenerator
	outputPrinter *format.OutputPrinter
	outWriter io.Writer
	errWriter io.Writer
}

func NewReqController(opts ReqControllerOptions) (obj *ReqController, err error) {
	obj = &ReqController{}

	// create a HTTP Invoker instance
	httpInvokerOptions := &engine.HttpInvokerOptions{}
	obj.httpInvoker, err = engine.NewHttpInvoker(httpInvokerOptions)
	if err != nil {
		return nil, err
	}

	// create a SpecGenerator instance
	obj.specGenerator, err = engine.NewSpecGenerator()
	if err != nil {
		return nil, err
	}
	if opts != nil {
		obj.specGenerator.Version = opts.GetVersion()
	}

	// create a OutputPrinter instance
	obj.outputPrinter, err = format.NewOutputPrinter(opts)
	if err != nil {
		return nil, err
	}

	return obj, err
}

func (r *ReqController) GetOutWriter() io.Writer {
	if r.outWriter == nil {
		return os.Stdout
	}
	return r.outWriter
}

func (r *ReqController) SetOutWriter(writer io.Writer) {
	r.outWriter = writer
}

func (r *ReqController) GetErrWriter() io.Writer {
	if r.errWriter == nil {
		return os.Stderr
	}
	return r.errWriter
}

func (r *ReqController) SetErrWriter(writer io.Writer) {
	r.errWriter = writer
}

func (z *ReqController) Execute(args ReqArguments) error {
	z.assertReady(args)

	generationPrinter := &GenerationPrinter{
		generator: z.specGenerator,
		writer: z.GetOutWriter(),
	}
	invocationPrinter := &InvocationPrinter{
		writer: z.GetOutWriter(),
	}

	if args.GetFormat() == "testcase" {
		_, err := z.httpInvoker.Do(transformReqArgs(args), generationPrinter)
		if err != nil {
			return z.displayError(err)
		}
		return nil
	}

	res, err := z.httpInvoker.Do(transformReqArgs(args), invocationPrinter)
	if err != nil {
		return z.displayError(err)
	}
	z.displayResult(res)
	return nil
}

type GenerationPrinter struct {
	generator *engine.SpecGenerator
	writer io.Writer
}

func (r *GenerationPrinter) PostProcess(req *engine.HttpRequest, res *engine.HttpResponse) error {
	if r.generator == nil {
		panic(fmt.Errorf("GenerationPrinter.generator must not be nil"))
	}
	if r.writer == nil {
		panic(fmt.Errorf("GenerationPrinter.writer must not be nil"))
	}
	return r.generator.GenerateTestCase(r.writer, req, res)
}

type InvocationPrinter struct {
	writer io.Writer
}

func (r *InvocationPrinter) PreProcess(req *engine.HttpRequest) error {
	return renderRequest(r.writer, req)
}

func (r *InvocationPrinter) PostProcess(req *engine.HttpRequest, res *engine.HttpResponse) error {
	return renderResponse(r.writer, res)
}

func renderRequest(w io.Writer, r *engine.HttpRequest) error {
	req, err := r.GetRawRequest()
	if err != nil {
		return err
	}
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

func renderResponse(w io.Writer, res *engine.HttpResponse) error {
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

func (z *ReqController) displayError(err error) error {
	if err == nil {
		return err
	}
	w := z.GetErrWriter()
	if w == nil {
		return err
	}
	fmt.Fprintf(w, "* %s\n", err.Error())
	return err
}

func (z *ReqController) displayResult(res *engine.HttpResponse) *engine.HttpResponse {
	if res == nil {
		return res
	}
	w := z.GetOutWriter()
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
