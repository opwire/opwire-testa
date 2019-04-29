package bootstrap

import (
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/format"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/utils"
)

type GenControllerOptions interface {
	script.Source
	GetNoColor() bool
	GetVersion() string
}

type GenController struct {
	scriptLoader *script.Loader
	scriptSelector *script.Selector
	scriptSource script.Source
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
	ref.scriptLoader, err = script.NewLoader(nil)
	if err != nil {
		return nil, err
	}

	// create a Script Selector instance
	ref.scriptSelector, err = script.NewSelector(opts)
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

func (c *GenController) Execute(args GenArguments) error {
	var testDirs []string
	if c.scriptSource != nil {
		testDirs = c.scriptSource.GetTestDirs()
	}

	var testFile string
	if args != nil {
		testFile = args.GetTestFile()
	}

	// display environment of command
	c.outputPrinter.Println()
	c.outputPrinter.Println(c.outputPrinter.Heading("Context"))

	relaDirs := utils.DetectRelativePaths(testDirs)
	if relaDirs != nil && len(relaDirs) > 0 {
		c.outputPrinter.Println(c.outputPrinter.ContextInfo("Test directories", "", relaDirs...))
	} else {
		c.outputPrinter.Println(c.outputPrinter.ContextInfo("Test directories", "Unspecified"))
	}

	if len(testFile) > 0 {
		c.outputPrinter.Println(c.outputPrinter.ContextInfo("File filter", testFile))
	}

	c.outputPrinter.Println(c.outputPrinter.ContextInfo("Name filter (" + c.scriptSelector.TypeOfTestNameFilter() + ")", c.scriptSelector.GetTestNameFilter()))

	// display prerequisites
	c.outputPrinter.Println()
	c.outputPrinter.Println(c.outputPrinter.Heading("Loading"))

	// Load testing script files from "test-dirs"
	descriptors := c.scriptLoader.LoadFrom(testDirs)

	// filter invalid descriptors and display errors
	descriptors, rejected := filterInvalidDescriptors(descriptors)

	for _, d := range rejected {
		c.outputPrinter.Println(c.outputPrinter.TestSuiteTitle(d.Locator.RelativePath))
		c.outputPrinter.Println(c.outputPrinter.Section(d.Error.Error()))
	}

	// filter testing script files by "test-file"
	descriptors = filterDescriptorsByFilePattern(descriptors, testFile)

	// filter target testcase by "test-name" title/name
	testcases := c.scriptSelector.GetTestCases(descriptors)

	// running & result
	c.outputPrinter.Println()
	c.outputPrinter.Println(c.outputPrinter.Heading("Running"))

	// raise an error if testcase not found or more than one found
	if len(testcases) == 0 {
		c.outputPrinter.Println(c.outputPrinter.ContextInfo("Error", "There is no testcase satisfied criteria"))
	}
	if len(testcases) > 1 {
		testinfo := make([]string, len(testcases))
		for i, test := range testcases {
			testinfo[i] = test.Title
			if len(test.Tags) > 0 {
				testinfo[i] += " (tags: " + strings.Join(test.Tags, ", ") + ")"
			}
		}
		c.outputPrinter.Println(c.outputPrinter.ContextInfo("Error", "There are more than one testcases satisfied criteria", testinfo...))
	}

	// generate curl statement from testcase's request
	if len(testcases) == 1 {
		testcase := testcases[0]
		request := testcase.Request

		generator := new(CurlGenerator)
		generator.generateCommand(os.Stdout, request)
	}

	c.outputPrinter.Println()
	return nil
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
