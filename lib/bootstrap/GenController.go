package bootstrap

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/script"
)

type GenControllerOptions interface {
	GetVersion() string
}

type GenController struct {
	loader *script.Loader
}

func NewGenController(opts GenControllerOptions) (ref *GenController, err error) {
	ref = &GenController{}

	ref.loader, err = script.NewLoader(nil)
	if err != nil {
		return nil, err
	}

	return ref, err
}

type GenArguments interface {
	GetTestDirs() []string
	GetTestFile() string
	GetTestCase() string
}

func (c *GenController) Execute(args GenArguments) error {
	// Load testing script files from "test-dirs"
	var testDirs []string
	if args != nil {
		testDirs = args.GetTestDirs()
	}
	descriptors := c.loader.LoadScripts(testDirs)

	// filter invalid descriptors and display errors
	descriptors = filterInvalidDescriptors(descriptors)

	// filter testing script files by "test-file"
	var testFile string
	if args != nil {
		testFile = args.GetTestFile()
	}
	if len(testFile) > 0 {
		descriptors = filterDescriptorsBySuffix(descriptors, testFile)
	}

	// filter target testcase by "test-case" title/name
	var testName string
	if args != nil {
		testName = args.GetTestCase()
	}
	if len(testName) > 0 {
		testName = standardizeName(testName)
	}

	testcases := make([]*engine.TestCase, 0)
	for _, d := range descriptors {
		testsuite := d.TestSuite
		if testsuite != nil {
			for _, testcase := range testsuite.TestCases {
				name := standardizeName(testcase.Title)
				if len(testName) == 0 || strings.HasPrefix(name, testName) {
					testcases = append(testcases, testcase)
				}
			}
		}
	}

	// raise an error if testcase not found or more than one found
	if len(testcases) == 0 {
		fmt.Printf("There are no testcases matched [%s]\n", testName)
		return nil
	}
	if len(testcases) > 1 {
		fmt.Printf("There are more than one testcase matched [%s]\n", testName)
		return nil
	}

	// generate curl statement from testcase's request
	testcase := testcases[0]
	request := testcase.Request

	generator := new(CurlGenerator)
	generator.generateCommand(os.Stdout, request)

	return nil
}

func filterInvalidDescriptors(src map[string]*script.Descriptor) map[string]*script.Descriptor {
	dst := make(map[string]*script.Descriptor, 0)
	for key, d := range src {
		if d.Error != nil {
			fmt.Printf("[%s] loading has been failed, error: %s\n", d.Locator.RelativePath, d.Error)
		} else {
			dst[key] = d
		}
	}
	return dst
}

func filterDescriptorsBySuffix(src map[string]*script.Descriptor, suffix string) map[string]*script.Descriptor {
	dst := make(map[string]*script.Descriptor, 0)
	for key, d := range src {
		if strings.HasSuffix(d.Locator.FullPath, suffix) {
			dst[key] = d
			continue
		}
		matched, err := filepath.Match(suffix, d.Locator.FullPath)
		if (err == nil && matched) {
			dst[key] = d
			continue
		}
	}
	return dst
}

func standardizeName(name string) string {
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
