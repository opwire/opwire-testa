package bootstrap

import (
	"flag"
	"fmt"
	"testing"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/script"
	"github.com/opwire/opwire-testa/lib/utils"
)

type RunControllerOptions interface {
	GetConfigPath() string
	GetVersion() string
	GetRevision() string
}

type RunController struct {
	options RunControllerOptions
	loader *script.Loader
	handler *engine.SpecHandler
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
	r.handler, err = engine.NewSpecHandler(r.storage)
	if err != nil {
		return nil, err
	}

	return r, nil
}

type RunArguments interface {
	GetSpecDirs() []string
}

func (r *RunController) Execute(args RunArguments) error {
	flag.Set("test.v", "true")

	var specDirs []string
	if args != nil {
		specDirs = args.GetSpecDirs()
	}

	// load test specifications
	allscripts := r.loader.LoadScripts(specDirs)

	// filter valid descriptors and display errors
	descriptors := make(map[string]*script.Descriptor, 0)
	for key, descriptor := range allscripts {
		if descriptor.Error != nil {
			rel, _ := utils.DetectRelativePath(descriptor.Locator.FullPath)
			fmt.Printf("Load [%s] has been failed, error: %s\n", rel, descriptor.Error)
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
	if r.handler == nil {
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
			result, err := r.handler.Examine(testcase)
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
