package bootstrap

import (
	"flag"
	"fmt"
	"testing"
	"github.com/opwire/opwire-qakit/lib/engine"
	"github.com/opwire/opwire-qakit/lib/script"
)

type TestRunnerOptions interface {
	GetSpecDirs() []string
}

type TestRunner struct {
	options TestRunnerOptions
	specDirs []string
	loader *script.Loader
	handler *engine.SpecHandler
	storage *TestStorage
}

func NewTestRunner(opts TestRunnerOptions) (r *TestRunner, err error) {
	r = &TestRunner{}

	// determine test specification dirs
	if opts == nil {
		r.specDirs = []string{}
	} else {
		r.specDirs = opts.GetSpecDirs()
	}

	// testing temporary storage
	r.storage = &TestStorage{}

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

func (r *TestRunner) loadTestSuites() (map[string]*script.Descriptor, error) {
	return r.loader.LoadScripts(r.specDirs)
}

func (r *TestRunner) wrapTestSuites(descriptors map[string]*script.Descriptor) ([]testing.InternalTest, error) {
	if r.handler == nil {
		return nil, fmt.Errorf("SpecHandler must not be nil")
	}
	tests := make([]testing.InternalTest, 0)
	for _, descriptor := range descriptors {
		for _, scenario := range descriptor.Scenarios {
			tests = append(tests, r.wrapTestCase(scenario))
		}
	}
	return tests, nil
}

func (r *TestRunner) wrapTestCase(scenario *engine.Scenario) (testing.InternalTest) {
	return testing.InternalTest{
		Name: scenario.Title,
		F: func (t *testing.T) {
			result, err := r.handler.Examine(scenario)
			_ = result
			_ = err
		},
	}
}

func (a *TestRunner) RunTests() error {
	flag.Set("test.v", "true")
	descriptors, err := a.loadTestSuites()
	if err != nil {
		return err
	}
	internalTests, err2 := a.wrapTestSuites(descriptors)
	if err2 != nil {
		return err2
	}
	testing.Main(defaultMatchString, internalTests, []testing.InternalBenchmark{}, []testing.InternalExample{})
	return nil
}

func defaultMatchString(pat, str string) (bool, error) {
	return true, nil
}

type TestStorage struct {}

