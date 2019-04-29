package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/format"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/tag"
	"github.com/opwire/opwire-testa/lib/utils"
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

type GenArguments interface {
	GetTestFile() string
}

func (r *GenController) Execute(args GenArguments) error {
	var testDirs []string
	if r.scriptSource != nil {
		testDirs = r.scriptSource.GetTestDirs()
	}

	var testFile string
	if args != nil {
		testFile = args.GetTestFile()
	}

	// display environment of command
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Context"))

	relaDirs := utils.DetectRelativePaths(testDirs)
	if relaDirs != nil && len(relaDirs) > 0 {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Test directories", "", relaDirs...))
	} else {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Test directories", "Unspecified"))
	}

	if len(testFile) > 0 {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("File filter", testFile))
	}

	r.outputPrinter.Println(r.outputPrinter.ContextInfo("Name filter (" + r.scriptSelector.TypeOfTestNameFilter() + ")", r.scriptSelector.GetTestNameFilter()))

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

	// filter testing script files by "test-file"
	descriptors = filterDescriptorsByFilePattern(descriptors, testFile)

	// filter target testcase by "test-name" title/name
	testcases := r.scriptSelector.GetTestCases(descriptors)

	// filter testcases by conditional tags
	testcases, _ = r.filterTestCasesByTags(testcases)

	// running & result
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Running"))

	// raise an error if testcase not found or more than one found
	if len(testcases) == 0 {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Error", "There is no testcase satisfied criteria"))
	}
	if len(testcases) > 1 {
		testinfo := make([]string, len(testcases))
		for i, test := range testcases {
			testinfo[i] = test.Title
			if len(test.Tags) > 0 {
				testinfo[i] += " (tags: " + strings.Join(test.Tags, ", ") + ")"
			}
		}
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Error", "There are more than one testcases satisfied criteria", testinfo...))
	}

	// generate curl statement from testcase's request
	if len(testcases) == 1 {
		testcase := testcases[0]
		request := testcase.Request

		generator := new(CurlGenerator)
		generator.generateCommand(os.Stdout, request)
	}

	r.outputPrinter.Println()
	return nil
}

func (r *GenController) filterTestCasesByTags(testcases []*engine.TestCase) (accepted []*engine.TestCase, rejected []*engine.TestCase) {
	accepted = make([]*engine.TestCase, 0)
	rejected = make([]*engine.TestCase, 0)
	for _, testcase := range testcases {
		if len(testcase.Tags) == 0 || r.tagManager.IsActive(testcase.Tags) {
			accepted = append(accepted, testcase)
		} else {
			rejected = append(rejected, testcase)
		}
	}
	return accepted, rejected
}

type CurlGenerator struct {
}

func (g *CurlGenerator) generateCommand(w io.Writer, req *engine.HttpRequest) error {
	fmt.Fprintf(w, "curl \\\n")
	fmt.Fprintf(w, "  --request %s \\\n", req.Method)
	fmt.Fprintf(w, "  --url \"%s\" \\\n", engine.BuildUrl(req, "", ""))
	for _, header := range req.Headers {
		fmt.Fprintf(w, "  --header '%s: %s' \\\n", header.Name, header.Value)
	}
	fmt.Fprintf(w, "  --data='%s'\n", req.Body)
	fmt.Fprintln(w)
	return nil
}
