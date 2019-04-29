package bootstrap

import (
	"flag"
	"fmt"
	"strings"
	"testing"
	"github.com/opwire/opwire-testa/lib/format"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/utils"
)

type RunControllerOptions interface {
	GetConfigPath() string
	GetConditionalTags() []string
	GetNoColor() bool
	GetVersion() string
	GetRevision() string
}

type RunController struct {
	loader *script.Loader
	specHandler *engine.SpecHandler
	tagManager *engine.TagManager
	outputPrinter *format.OutputPrinter
	storage *TestStateStore
}

func NewRunController(opts RunControllerOptions) (r *RunController, err error) {
	r = &RunController{}

	// testing temporary storage
	r.storage = &TestStateStore{}

	// create a Script Loader instance
	r.loader, err = script.NewLoader(r.storage)
	if err != nil {
		return nil, err
	}

	// create a Spec Handler instance
	r.specHandler, err = engine.NewSpecHandler(r.storage)
	if err != nil {
		return nil, err
	}

	// create a TagManager instance
	r.tagManager, err = engine.NewTagManager(opts)
	if err != nil {
		return nil, err
	}

	// create a OutputPrinter instance
	r.outputPrinter, err = format.NewOutputPrinter(opts)
	if err != nil {
		return nil, err
	}

	return r, nil
}

type RunArguments interface {
	GetTestDirs() []string
	GetTestName() string
}

func (r *RunController) Execute(args RunArguments) error {
	flag.Set("test.v", "false")

	// begin environments
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Context"))

	var testDirs []string
	if args != nil {
		testDirs = args.GetTestDirs()
	}
	relaDirs := utils.DetectRelativePaths(testDirs)
	if relaDirs != nil && len(relaDirs) > 0 {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Test directories", "", relaDirs...))
	} else {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Test directories", "Unspecified"))
	}

	inclTags := r.tagManager.GetIncludedTags()
	if inclTags != nil && len(inclTags) > 0 {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Selected tags", strings.Join(inclTags, ", ")))
	} else {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Selected tags", "Unspecified"))
	}

	exclTags := r.tagManager.GetExcludedTags()
	if exclTags != nil && len(exclTags) > 0 {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Excluded tags", strings.Join(exclTags, ", ")))
	} else {
		r.outputPrinter.Println(r.outputPrinter.ContextInfo("Excluded tags", "Unspecified"))
	}

	scriptSelector, err := script.NewSelector(args)
	if err != nil {
		return err
	}
	_ = scriptSelector

	// begin prerequisites
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Loading"))

	// load test specifications
	allscripts := r.loader.LoadScripts(testDirs)

	// filter invalid descriptors and display errors
	descriptors := make(map[string]*script.Descriptor, 0)
	for key, descriptor := range allscripts {
		if descriptor.Error != nil {
			r.outputPrinter.Println(r.outputPrinter.TestSuiteTitle(descriptor.Locator.RelativePath))
			r.outputPrinter.Println(r.outputPrinter.Section(descriptor.Error.Error()))
		} else {
			descriptors[key] = descriptor
		}
	}

	// begin testing
	r.outputPrinter.Println()
	r.outputPrinter.Println(r.outputPrinter.Heading("Testing"))

	// create the test runners
	internalTests, err2 := r.wrapTestSuites(descriptors)
	if err2 != nil {
		return err2
	}

	// Run the tests
	testing.Main(defaultMatchString, internalTests, []testing.InternalBenchmark{}, []testing.InternalExample{})

	// endof testing
	r.outputPrinter.Println()

	return nil
}

func defaultMatchString(pat, str string) (bool, error) {
	return true, nil
}

type TestStateStore struct {}

func (r *RunController) wrapTestSuites(descriptors map[string]*script.Descriptor) ([]testing.InternalTest, error) {
	if r.specHandler == nil {
		panic(fmt.Errorf("SpecHandler must not be nil"))
	}
	tests := make([]testing.InternalTest, 0)
	for _, descriptor := range descriptors {
		test, err := r.wrapDescriptor(descriptor)
		if err == nil {
			tests = append(tests, test)
		}
	}
	return tests, nil
}

func (r *RunController) wrapDescriptor(descriptor *script.Descriptor) (testing.InternalTest, error) {
	testsuite := descriptor.TestSuite
	if testsuite == nil {
		return testing.InternalTest{}, fmt.Errorf("TestSuite must not be nil")
	}
	return testing.InternalTest{
		Name: descriptor.Locator.RelativePath,
		F: func (t *testing.T) {
			r.outputPrinter.Println(r.outputPrinter.TestSuiteTitle(descriptor.Locator.RelativePath))
			tests := make([]testing.InternalTest, 0)
			for _, testcase := range testsuite.TestCases {
				tests = append(tests, r.wrapTestCase(testcase))
			}
			testing.RunTests(defaultMatchString, tests)
			r.outputPrinter.Println()
		},
	}, nil
}

func (r *RunController) wrapTestCase(testcase *engine.TestCase) (testing.InternalTest) {
	return testing.InternalTest{
		Name: testcase.Title,
		F: func (t *testing.T) {
			ok := true
			if len(testcase.Tags) > 0 && !r.tagManager.IsActive(testcase.Tags) {
				ok = false
				r.outputPrinter.Println(r.outputPrinter.Skipped(testcase.Title))
			}
			result, err := r.specHandler.Examine(testcase)
			if err != nil {
				ok = false
				r.outputPrinter.Println(r.outputPrinter.Cracked(testcase.Title))
			}
			if result != nil && len(result.Errors) > 0 {
				ok = false
				r.outputPrinter.Println(r.outputPrinter.Failure(testcase.Title))
				for key, err := range result.Errors {
					r.outputPrinter.Printf(r.outputPrinter.SectionTitle(key))
					r.outputPrinter.Printf(r.outputPrinter.Section(err.Error()))
				}
			}
			if ok {
				r.outputPrinter.Println(r.outputPrinter.Success(testcase.Title))
			}
		},
	}
}
