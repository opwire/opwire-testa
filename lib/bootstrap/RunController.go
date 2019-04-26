package bootstrap

import (
	"flag"
	"fmt"
	"testing"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/script"
)

type RunControllerOptions interface {
	GetConfigPath() string
	GetConditionalTags() []string
	GetVersion() string
	GetRevision() string
}

type RunController struct {
	loader *script.Loader
	specHandler *engine.SpecHandler
	tagManager *engine.TagManager
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

	return r, nil
}

type RunArguments interface {
	GetTestDirs() []string
}

func (r *RunController) Execute(args RunArguments) error {
	flag.Set("test.v", "true")

	var testDirs []string
	if args != nil {
		testDirs = args.GetTestDirs()
	}

	// load test specifications
	allscripts := r.loader.LoadScripts(testDirs)

	// filter invalid descriptors and display errors
	descriptors := make(map[string]*script.Descriptor, 0)
	for key, descriptor := range allscripts {
		if descriptor.Error != nil {
			fmt.Printf("[%s] loading has been failed, error: %s\n", descriptor.Locator.RelativePath, descriptor.Error)
		} else {
			descriptors[key] = descriptor
		}
	}

	// create the test runners
	internalTests, err2 := r.wrapTestSuites(descriptors)
	if err2 != nil {
		return err2
	}

	// Run the tests
	testing.Main(defaultMatchString, internalTests, []testing.InternalBenchmark{}, []testing.InternalExample{})
	return nil
}

func defaultMatchString(pat, str string) (bool, error) {
	return true, nil
}

type TestStateStore struct {}

func (r *RunController) wrapTestSuites(descriptors map[string]*script.Descriptor) ([]testing.InternalTest, error) {
	if r.specHandler == nil {
		return nil, fmt.Errorf("SpecHandler must not be nil")
	}
	tests := make([]testing.InternalTest, 0)
	for _, descriptor := range descriptors {
		testsuite := descriptor.TestSuite
		if testsuite != nil {
			for _, testcase := range testsuite.TestCases {
				tests = append(tests, r.wrapTestCase(testcase))
			}
		}
	}
	return tests, nil
}

func (r *RunController) wrapTestCase(testcase *engine.TestCase) (testing.InternalTest) {
	return testing.InternalTest{
		Name: testcase.Title,
		F: func (t *testing.T) {
			if len(testcase.Tags) > 0 {
				if !r.tagManager.IsActive(testcase.Tags) {
					t.Skipf("[%s] is disabled by conditional tags", testcase.Title)
				}
			}
			result, err := r.specHandler.Examine(testcase)
			if err != nil {
				t.Fatalf("[%s] got a fatal error. Exit now", testcase.Title)
			}
			if result != nil && len(result.Errors) > 0 {
				t.Errorf("[%s] has failed:", testcase.Title)
				for key, err := range result.Errors {
					t.Logf("+ %s", key)
					t.Logf("|- %s", err)
				}
			}
		},
	}
}
