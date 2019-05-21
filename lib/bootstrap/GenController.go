package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/opwire/opwire-testa/lib/client"
	"github.com/opwire/opwire-testa/lib/format"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/tag"
)

type GenControllerOptions interface {
	script.Source
	GetNoColor() bool
}

type GenController struct {
	scriptLoader *script.Loader
	scriptSelector *script.Selector
	scriptSource script.Source
	tagManager *tag.Manager
	outputPrinter *format.OutputPrinter
	outWriter io.Writer
}

func NewGenController(opts GenControllerOptions) (ref *GenController, err error) {
	ref = &GenController{}

	// testing temporary storage
	ref.scriptSource, err = script.NewSource(opts)
	if err != nil {
		return nil, err
	}

	// create a Script Loader instance
	ref.scriptLoader, err = script.NewLoader(ref.scriptSource)
	if err != nil {
		return nil, err
	}

	// create a Script Selector instance
	ref.scriptSelector, err = script.NewSelector(ref.scriptSource)
	if err != nil {
		return nil, err
	}

	// create a Manager instance
	ref.tagManager, err = tag.NewManager(ref.scriptSource)
	if err != nil {
		return nil, err
	}

	// create a OutputPrinter instance
	ref.outputPrinter, err = format.NewOutputPrinter(opts)
	if err != nil {
		return nil, err
	}

	return ref, err
}

type GenArguments interface {}

func (r *GenController) GetOutWriter() io.Writer {
	if r.outWriter == nil {
		return os.Stdout
	}
	return r.outWriter
}

func (r *GenController) SetOutWriter(writer io.Writer) {
	r.outWriter = writer
}

func (r *GenController) Execute(args GenArguments) error {
	// display environment of command
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Context"))
	printScriptSourceArgs(r.outputPrinter, r.scriptSource, r.scriptSelector, r.tagManager)

	// display prerequisites
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Loading"))

	// Load testing script files from "test-dirs"
	descriptors := r.scriptLoader.Load()

	// filter invalid descriptors and display errors
	descriptors, rejected := filterInvalidDescriptors(descriptors)
	for _, d := range rejected {
		r.outputPrinter.Println(r.outputPrinter.TestSuiteTitle(d.Locator.RelativePath))
		r.outputPrinter.Println(r.outputPrinter.Section(d.Error.Error()))
	}

	// filter testing script files by "inclusive-files"
	descriptors = filterDescriptorsByInclusivePatterns(descriptors, r.scriptSource.GetInclFiles())

	// filter testing script files by "exclusive-files"
	descriptors = filterDescriptorsByExclusivePatterns(descriptors, r.scriptSource.GetExclFiles())

	// filter target testcase by "test-name" title/name
	testcases := r.scriptSelector.GetTestCases(descriptors)

	// filter testcases by conditional tags
	testcases, _ = filterTestCasesByTags(r.tagManager, testcases)

	// running & result
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Running"))

	// raise an error if testcase not found or more than one found
	if len(testcases) == 0 {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Error", "There is no testcase satisfied criteria"))
	}
	if len(testcases) > 1 {
		testinfo := make([]string, len(testcases))
		for i, testcase := range testcases {
			testinfo[i] = testcase.Title
			if len(testcase.Tags) > 0 {
				testinfo[i] += " (tags: " + strings.Join(testcase.Tags, ", ") + ")"
			}
		}
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Error", "There are more than one testcases satisfied criteria", testinfo...))
	}

	// generate curl statement from testcase's request
	if len(testcases) == 1 {
		testcase := testcases[0]
		request := testcase.Request

		generator := new(CurlGenerator)
		generator.generateCommand(r.GetOutWriter(), request)
	}

	r.outputPrinter.Println()
	return nil
}

type CurlGenerator struct {}

func (g *CurlGenerator) generateCommand(w io.Writer, req *client.HttpRequest) error {
	fmt.Fprintf(w, "curl \\\n")
	fmt.Fprintf(w, "  --request %s \\\n", req.Method)
	fmt.Fprintf(w, "  --url \"%s\" \\\n", client.BuildUrl(req))
	for _, header := range req.Headers {
		fmt.Fprintf(w, "  --header '%s: %s' \\\n", header.Name, header.Value)
	}
	fmt.Fprintf(w, "  --data='%s'\n", req.Body)
	return nil
}
