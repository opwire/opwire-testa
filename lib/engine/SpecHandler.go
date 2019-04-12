package engine

import(
	"fmt"
)

type SpecHandlerOptions interface {
}

type SpecHandler struct {
	invoker *HttpInvoker
}

func NewSpecHandler(opts SpecHandlerOptions) (e *SpecHandler, err error) {
	e = &SpecHandler{}
	e.invoker, err = NewHttpInvoker(nil)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *SpecHandler) Examine(scenario *Scenario) (*ExaminationResult, error) {
	if scenario == nil {
		return nil, fmt.Errorf("Scenario must not be nil")
	}
	result := &ExaminationResult{}
	res, err := e.invoker.Do(scenario.Request)
	if err != nil {
		return nil, err
	}
	result.Response = res
	return result, nil
}

type Scenario struct {
	Title string `yaml:"title"`
	Skipped bool `yaml:"skipped"`
	Request *HttpRequest `yaml:"request"`
	Measure *HttpMeasure `yaml:"measure"`
}

type ExaminationResult struct {
	Response *HttpResponse
}
