package engine

import(
	"fmt"
	"testing"
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

func (e *SpecHandler) Examine(t *testing.T, scenario *Scenario) (*ExaminationResult, error) {
	if scenario == nil {
		return nil, fmt.Errorf("Scenario must not be nil")
	}

	result := &ExaminationResult{}

	// make the testing request
	res, err := e.invoker.Do(scenario.Request)
	if err != nil {
		return nil, err
	}
	result.Response = res

	// matching with expectation
	return result, nil
}

type Scenario struct {
	Title string `yaml:"title"`
	Skipped bool `yaml:"skipped"`
	Request *HttpRequest `yaml:"request"`
	Expectation *Expectation `yaml:"expectation"`
}

type Expectation struct {
	StatusCode *MeasureStatusCode `yaml:"status-code"`
	Headers *MeasureHeaders `yaml:"headers"`
	Body *MeasureBody `yaml:"body"`
}

type MeasureStatusCode struct {
	EqualTo *int `yaml:"equal-to"`
}

type MeasureHeaders struct {
	HasTotal *int `yaml:"has-total"`
	Items []MeasureHeader `yaml:"items"`
}

type MeasureHeader struct {
	Name *string `yaml:"name"`
	EqualTo *string `yaml:"equal-to"`
}

type MeasureBody struct {
	HasFormat *string `yaml:"has-format"`
	EqualTo *string `yaml:"equal-to"`
	MatchWith *string `yaml:"match-with"`
}

type ExaminationResult struct {
	Response *HttpResponse
}
