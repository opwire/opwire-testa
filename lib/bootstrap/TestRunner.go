package bootstrap

import (
	"flag"
	"fmt"
	"testing"
	"github.com/opwire/opwire-qakit/lib/engine"
	"github.com/opwire/opwire-qakit/lib/script"
)

type TestRunnerOptions interface {}

type TestRunner struct {
	handler *engine.SpecHandler
}

func NewTestRunner(opts TestRunnerOptions) (*TestRunner, error) {
	r := &TestRunner{}
	return r, nil
}

func (r *TestRunner) Export(descriptors map[string]*script.Descriptor) ([]testing.InternalTest, error) {
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
	tests, err := a.Export(nil)
	if err != nil {
		return err
	}
	testing.Main(defaultMatchString, tests, []testing.InternalBenchmark{}, []testing.InternalExample{})
	return nil
}

func defaultMatchString(pat, str string) (bool, error) {
	return true, nil
}

func Test1(t *testing.T) {
	if 1+2 != 3 {
		t.Fail()
	}
}

func Test2(t *testing.T) {
	if 3*3 == 9 {
		t.Fail() // WHOOPS!
	}
}

func Test3(t *testing.T) {
	fmt.Println("Just wanted to print here")
}
