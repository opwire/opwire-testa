package engine

import(
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
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

	// make the testing request
	res, err := e.invoker.Do(scenario.Request)
	if err != nil {
		return nil, err
	}
	result.Response = res

	// matching with expectation
	errors := make(map[string]error, 0)
	expect := scenario.Expectation
	if expect != nil {
		_sc := expect.StatusCode
		if _sc != nil {
			if _sc.EqualTo != nil {
				if res.StatusCode != *_sc.EqualTo {
					errors["StatusCode"] = fmt.Errorf("Response StatusCode [%d] is not equal to expected value [%d]", res.StatusCode, *_sc.EqualTo)
				}
			}
		}
		_hs := expect.Headers
		if _hs != nil {
			if _hs.Items != nil {
				for _, item := range _hs.Items {
					headerVal := res.Header.Get(*item.Name)
					if item.EqualTo != nil && *item.EqualTo != headerVal {
						errors[fmt.Sprintf("Header[%s]", *item.Name)] = fmt.Errorf("Returned value: [%s] is mismatched with expected: [%s]", headerVal, *item.EqualTo)
					}
				}
			}
		}
		_eb := expect.Body
		if _eb != nil {
			if _eb.EqualTo != nil {
				_rb := string(res.Body)
				if res.Body == nil || _rb != *_eb.EqualTo {
					errors["Body"] = fmt.Errorf("Response body is mismatched with expected content.\n    Received: %s\n    Expected: %s", _rb, *_eb.EqualTo)
				}
			}
			if _eb.JSONEquals != nil {
				var receivedJSON, expectedJSON interface{}
				next := true
				if (res.Body == nil) {
					errors["Body/JSONEquals/empty"] = fmt.Errorf("Response body is empty (invalid JSON)")
					next = false
				}
				if err := json.Unmarshal(res.Body, &receivedJSON); err != nil {
					errors["Body/JSONEquals/received"] = fmt.Errorf("Invalid response content: %s", err)
					next = false
				}
				if err := json.Unmarshal([]byte(*_eb.JSONEquals), &expectedJSON); err != nil {
					errors["Body/JSONEquals/expected"] = fmt.Errorf("Invalid expected content: %s", err)
					next = false
				}
				if diff := cmp.Diff(expectedJSON, receivedJSON); diff != "" {
					errors["Body/JSONEquals/result"] = fmt.Errorf("Body mismatch (-expected +received):\n%s", diff)
				}
				_ = next
			}
		}
	}
	result.Errors = errors

	return result, nil
}

type Scenario struct {
	Title string `yaml:"title"`
	Skipped bool `yaml:"skipped"`
	OnError string `yaml:"on-error"`
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
	MatchWith *string `yaml:"match-with"`
	EqualTo *string `yaml:"equal-to"`
	JSONEquals *string `yaml:"equal-json"`
	JSONContains *string `yaml:"json-contains"`
	YAMLEquals *string `yaml:"equal-yaml"`
	YAMLContains *string `yaml:"yaml-contains"`
}

type ExaminationResult struct {
	Errors map[string]error
	Response *HttpResponse
}
