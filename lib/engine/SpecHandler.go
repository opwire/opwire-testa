package engine

import(
	"encoding/json"
	"fmt"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/google/go-cmp/cmp"
	"github.com/opwire/opwire-testa/lib/comparison"
	"github.com/opwire/opwire-testa/lib/utils"
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

func (e *SpecHandler) Examine(testcase *TestCase) (*ExaminationResult, error) {
	if testcase == nil {
		return nil, fmt.Errorf("TestCase must not be nil")
	}

	result := &ExaminationResult{}

	// check if testcase is pending
	if testcase.Pending != nil && *testcase.Pending == true {
		result.Status = "pending"
		return result, nil
	}

	// start time
	startTime := time.Now()

	// make the testing request
	res, err := e.invoker.Do(testcase.Request)
	if err != nil {
		result.Duration = time.Since(startTime)
		return nil, err
	}
	result.Response = res

	// matching with expectation
	errors := make(map[string]error, 0)
	expect := testcase.Expectation
	if expect != nil {
		_sc := expect.StatusCode
		if _sc != nil && _sc.Is != nil {
			if _sc.Is.EqualTo != nil {
				if eq, _ := comparison.IsEqualTo(res.StatusCode, _sc.Is.EqualTo); !eq {
					errors["StatusCode"] = fmt.Errorf("Response StatusCode [%d] is not equal to expected value [%v]", res.StatusCode, _sc.Is.EqualTo)
				}
			}
			if _sc.Is.ContainedIn != nil {
				if !comparison.BelongsTo(res.StatusCode, _sc.Is.ContainedIn) {
					errors["StatusCode"] = fmt.Errorf("Response StatusCode [%d] does not belong to expected list %v", res.StatusCode, _sc.Is.ContainedIn)
				}
			}
		}
		_hs := expect.Headers
		if _hs != nil {
			if _hs.Total != nil && _hs.Total.Is != nil {
				headerTotal := len(res.Header)
				totalIs := _hs.Total.Is
				if totalIs.EqualTo != nil {
					if eq, _ := comparison.IsEqualTo(headerTotal, totalIs.EqualTo); !eq {
						errors["Header/Total"] = fmt.Errorf("Total of headers (%d) mismatchs with expected number (%v)", headerTotal, totalIs.EqualTo)
					}
				}
			}
			if _hs.Items != nil {
				for _, item := range _hs.Items {
					headerVal := res.Header.Get(*item.Name)
					if item.Is != nil && item.Is.EqualTo != nil {
						eq, _ := comparison.IsEqualTo(headerVal, item.Is.EqualTo)
						if !eq {
							errors[fmt.Sprintf("Header[%s]", *item.Name)] = fmt.Errorf("Returned value: [%s] is mismatched with expected: [%s]", headerVal, item.Is.EqualTo)
						}
					}
				}
			}
		}
		_eb := expect.Body
		if _eb != nil && _eb.HasFormat != nil {
			var format string = *_eb.HasFormat
			if format == utils.BODY_FORMAT_FLAT {
				if _eb.IsEqualTo != nil {
					_rb := string(res.Body)
					if res.Body == nil || _rb != *_eb.IsEqualTo {
						errors["Body"] = fmt.Errorf("Response body is mismatched with expected content.\n    Received: %s\n    Expected: %s", _rb, *_eb.IsEqualTo)
					}
				}
			}
			if format == utils.BODY_FORMAT_JSON || format == utils.BODY_FORMAT_YAML {
				var receivedObj, expectedObj map[string]interface{}
				next := true
				if (res.Body == nil) {
					errors["Body/ReceivedObject"] = fmt.Errorf("Response body is empty (invalid %s)", format)
					next = false
				} else if err := Unmarshal(format, res.Body, &receivedObj); err != nil {
					errors["Body/ReceivedObject"] = fmt.Errorf("Invalid response content: %s", err)
					next = false
				}
				if next && _eb.IsEqualTo != nil {
					if err := Unmarshal(format, []byte(*_eb.IsEqualTo), &expectedObj); err != nil {
						errors["Body/ExpectedObject"] = fmt.Errorf("Invalid expected content: %s", err)
						next = false
					}
					if next {
						if diff := cmp.Diff(expectedObj, receivedObj); diff != "" {
							errors["Body/IsEqualTo"] = fmt.Errorf("Body mismatch (-expected +received):\n%s", diff)
						}
					}
				}
				if next && _eb.Includes != nil {
					if err := Unmarshal(format, []byte(*_eb.Includes), &expectedObj); err != nil {
						errors["Body/ReceivedObject"] = fmt.Errorf("Invalid expected content: %s", err)
						next = false
					}
					if next {
						ok, diff := comparison.IsPartOf(expectedObj, receivedObj)
						if !ok {
							errors["Body/Includes"] = fmt.Errorf("Body mismatch (-expected +received):\n%s", diff)
						}
					}
				}
				if next && len(_eb.Fields) > 0 {
					eFields := _eb.Fields
					rFields, _ := utils.Flatten("", receivedObj)
					for _, eField := range eFields {
						if eField.Is != nil && eField.Is.EqualTo != nil {
							eValue := eField.Is.EqualTo
							if rValue, ok := rFields[*eField.Path]; ok {
								if eq, _ := comparison.IsEqualTo(rValue, eValue); !eq {
									errors["Body/Fields/" + *eField.Path] = fmt.Errorf("Field mismatch expected: %v / received: %v", eValue, rValue)
								}
							} else {
								errors["Body/Fields/" + *eField.Path] = fmt.Errorf("Field not found, expected: %v", eValue)
							}
						}
					}
				}
			}
		} else {
			if _eb.HasFormat == nil && (_eb.IsEqualTo != nil || _eb.Includes != nil) {
				errors["Body/Expectation"] = fmt.Errorf("Unknown body format, please provides [has-format] value")
			}
		}
	}
	result.Errors = errors

	if len(errors) == 0 {
		result.Status = "ok"
	} else {
		result.Status = "error"
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

func Unmarshal(format string, source []byte, target interface{}) error {
	if format == utils.BODY_FORMAT_JSON {
		return json.Unmarshal(source, target)
	}
	if format == utils.BODY_FORMAT_YAML {
		return yaml.Unmarshal(source, target)
	}
	return fmt.Errorf("Invalid body format: %s", format)
}

type TestSuite struct {
	TestCases []*TestCase `yaml:"testcases" json:"testcases"`
	Pending *bool `yaml:"pending,omitempty" json:"pending"`
}

type TestCase struct {
	Title string `yaml:"title" json:"title"`
	Version *string `yaml:"version,omitempty" json:"version"`
	Request *HttpRequest `yaml:"request" json:"request"`
	Expectation *Expectation `yaml:"expectation" json:"expectation"`
	Pending *bool `yaml:"pending,omitempty" json:"pending"`
	Tags []string `yaml:"tags,omitempty" json:"tags"`
	CreatedTime *string `yaml:"created-time,omitempty" json:"created-time"`
}

type Expectation struct {
	StatusCode *MeasureStatusCode `yaml:"status-code,omitempty" json:"status-code"`
	Headers *MeasureHeaders `yaml:"headers,omitempty" json:"headers"`
	Body *MeasureBody `yaml:"body,omitempty" json:"body"`
}

type MeasureStatusCode struct {
	Is *ComparisonOperators `yaml:"is,omitempty" json:"is"`
}

type MeasureHeaders struct {
	Total *MeasureTotal `yaml:"total,omitempty" json:"total"`
	Items []MeasureHeader `yaml:"items,omitempty" json:"items"`
}

type MeasureTotal struct {
	Is *ComparisonOperators `yaml:"is,omitempty" json:"is"`
}

type MeasureHeader struct {
	Name *string `yaml:"name" json:"name"`
	Is *ComparisonOperators `yaml:"is,omitempty" json:"is"`
}

type MeasureBody struct {
	HasFormat *string `yaml:"has-format,omitempty" json:"has-format"`
	Includes *string `yaml:"includes,omitempty" json:"includes"`
	IsEqualTo *string `yaml:"is-equal-to,omitempty" json:"is-equal-to"`
	MatchWith *string `yaml:"match-with,omitempty" json:"match-with"`
	Fields []MeasureBodyField `yaml:"fields,omitempty" json:"fields"`
}

type MeasureBodyField struct {
	Path *string `yaml:"path,omitempty" json:"path"`
	Is *ComparisonOperators `yaml:"is,omitempty" json:"is"`
}

type ComparisonOperators struct {
	EqualTo interface{} `yaml:"equal-to,omitempty" json:"equal-to"`
	NotEqualTo interface{} `yaml:"not-equal-to,omitempty" json:"not-equal-to"`
	LT interface{} `yaml:"lt,omitempty" json:"lt"`
	LTE interface{} `yaml:"lte,omitempty" json:"lte"`
	GT interface{} `yaml:"gt,omitempty" json:"gt"`
	GTE interface{} `yaml:"gte,omitempty" json:"gte"`
	ContainedIn []interface{} `yaml:"contained-in,omitempty" json:"contained-in"`
}

type ExaminationResult struct {
	Duration time.Duration
	Errors map[string]error
	Response *HttpResponse
	Status string
}
