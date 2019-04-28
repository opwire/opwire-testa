package bootstrap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/format"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/utils"
)

type GenControllerOptions interface {
	GetNoColor() bool
	GetVersion() string
}

type GenController struct {
	loader *script.Loader
	outputPrinter *format.OutputPrinter
}

func NewGenController(opts GenControllerOptions) (ref *GenController, err error) {
	ref = &GenController{}

	// create a Script Loader instance
	ref.loader, err = script.NewLoader(nil)
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
	GetTestDirs() []string
	GetTestFile() string
	GetTestName() string
}

func (c *GenController) Execute(args GenArguments) error {
	var testDirs []string
	if args != nil {
		testDirs = args.GetTestDirs()
	}

	var testFile string
	if args != nil {
		testFile = args.GetTestFile()
	}

	var testName string
	if args != nil {
		testName = args.GetTestName()
	}

	var testNameRe *regexp.Regexp
	if utils.TEST_CASE_TITLE_REGEXP.MatchString(testName) {
		testName = standardizeName(testName)
	} else {
		re, err := regexp.Compile(testName)
		if err == nil {
			testNameRe = re
		}
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

	if len(testName) > 0 {
		c.outputPrinter.Println(c.outputPrinter.ContextInfo("Name filter", testName))
	}

	// display prerequisites
	c.outputPrinter.Println()
	c.outputPrinter.Println(c.outputPrinter.Heading("Loading"))

	// Load testing script files from "test-dirs"
	descriptors := c.loader.LoadScripts(testDirs)

	// filter invalid descriptors and display errors
	descriptors = c.filterInvalidDescriptors(descriptors)

	// filter testing script files by "test-file"
	if len(testFile) > 0 {
		descriptors = filterDescriptorsBySuffix(descriptors, testFile)
	}

	// filter target testcase by "test-name" title/name
	testcases := make([]*engine.TestCase, 0)
	for _, d := range descriptors {
		testsuite := d.TestSuite
		if testsuite != nil {
			for _, testcase := range testsuite.TestCases {
				if len(testName) == 0 {
					testcases = append(testcases, testcase)
				}
				if testNameRe == nil {
					name := standardizeName(testcase.Title)
					if strings.Contains(name, testName) {
						testcases = append(testcases, testcase)
					}
				} else {
					if testNameRe.MatchString(testcase.Title) {
						testcases = append(testcases, testcase)
					}
				}
			}
		}
	}

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

func (c *GenController) filterInvalidDescriptors(src map[string]*script.Descriptor) map[string]*script.Descriptor {
	dst := make(map[string]*script.Descriptor, 0)
	for key, d := range src {
		if d.Error != nil {
			c.outputPrinter.Println(c.outputPrinter.TestSuiteTitle(d.Locator.RelativePath))
			c.outputPrinter.Println(c.outputPrinter.Section(d.Error.Error()))
		} else {
			dst[key] = d
		}
	}
	return dst
}

func filterDescriptorsBySuffix(src map[string]*script.Descriptor, suffix string) map[string]*script.Descriptor {
	dst := make(map[string]*script.Descriptor, 0)
	for key, d := range src {
		if strings.HasSuffix(d.Locator.AbsolutePath, suffix) {
			dst[key] = d
			continue
		}
		matched, err := filepath.Match(suffix, d.Locator.AbsolutePath)
		if (err == nil && matched) {
			dst[key] = d
			continue
		}
	}
	return dst
}

func standardizeName(name string) string {
	if len(name) == 0 {
		return name
	}
	name = strings.ToLower(name)
	name = strings.Join(strings.Fields(strings.TrimSpace(name)), " ")
	return name
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
