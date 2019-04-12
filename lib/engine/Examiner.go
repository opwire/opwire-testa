package engine

type ExaminerOptions interface {
}

type Examiner struct {
	invoker *HttpInvoker
}

func NewExaminer(opts ExaminerOptions) (e *Examiner, err error) {
	e = &Examiner{}
	e.invoker, err = NewHttpInvoker(nil)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *Examiner) Examine(scenario Scenario) (*ExaminationResult, error) {
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
