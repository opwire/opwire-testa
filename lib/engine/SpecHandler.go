package engine

import(
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/google/go-cmp/cmp"
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
		if _sc != nil {
			if _sc.IsEqualTo != nil {
				if res.StatusCode != *_sc.IsEqualTo {
					errors["StatusCode"] = fmt.Errorf("Response StatusCode [%d] is not equal to expected value [%d]", res.StatusCode, *_sc.IsEqualTo)
				}
			}
			if _sc.BelongsTo != nil {
				if !utils.ContainsInt(_sc.BelongsTo, res.StatusCode) {
					errors["StatusCode"] = fmt.Errorf("Response StatusCode [%d] does not belong to expected list %v", res.StatusCode, _sc.BelongsTo)
				}
			}
		}
		_hs := expect.Headers
		if _hs != nil {
			if _hs.Items != nil {
				for _, item := range _hs.Items {
					headerVal := res.Header.Get(*item.Name)
					if item.IsEqualTo != nil && *item.IsEqualTo != headerVal {
						errors[fmt.Sprintf("Header[%s]", *item.Name)] = fmt.Errorf("Returned value: [%s] is mismatched with expected: [%s]", headerVal, *item.IsEqualTo)
					}
				}
			}
		}
		_eb := expect.Body
		if _eb != nil {
			if *_eb.HasFormat == "flat" {
				if _eb.IsEqualTo != nil {
					_rb := string(res.Body)
					if res.Body == nil || _rb != *_eb.IsEqualTo {
						errors["Body"] = fmt.Errorf("Response body is mismatched with expected content.\n    Received: %s\n    Expected: %s", _rb, *_eb.IsEqualTo)
					}
				}
			}
			if *_eb.HasFormat == "json" || *_eb.HasFormat == "yaml" {
				var format string = *_eb.HasFormat
				var receivedObj, expectedObj map[string]interface{}
				next := true
				if (res.Body == nil) {
					errors["Body/receivedObj"] = fmt.Errorf("Response body is empty (invalid %s)", format)
					next = false
				} else if err := Unmarshal(format, res.Body, &receivedObj); err != nil {
					errors["Body/receivedObj"] = fmt.Errorf("Invalid response content: %s", err)
					next = false
				}
				if next && _eb.IsEqualTo != nil {
					if err := Unmarshal(format, []byte(*_eb.IsEqualTo), &expectedObj); err != nil {
						errors["Body/expectedObj"] = fmt.Errorf("Invalid expected content: %s", err)
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
						errors["Body/expectedObj"] = fmt.Errorf("Invalid expected content: %s", err)
						next = false
					}
					if next {
						var r DiffReporter
						diff := cmp.Diff(expectedObj, receivedObj, cmp.Reporter(&r))
						if r.HasDiffs() {
							errors["Body/Includes"] = fmt.Errorf("Body mismatch (-expected +received):\n%s", diff)
						}
					}
				}
				if next && len(_eb.Fields) > 0 {
					eFields := _eb.Fields
					rFields, _ := utils.Flatten("", receivedObj)
					for _, eField := range eFields {
						eValue := eField.Value
						if rValue, ok := rFields[*eField.Path]; ok {
							if !IsEqual(rValue, eValue) {
								errors["Body/Fields/" + *eField.Path] = fmt.Errorf("Field mismatch expected: %v / received: %v", eValue, rValue)
							}
						} else {
							errors["Body/Fields/" + *eField.Path] = fmt.Errorf("Field not found, expected: %v", eValue)
						}
					}
				}
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
	if format == "json" {
		return json.Unmarshal(source, target)
	}
	if format == "yaml" {
		return yaml.Unmarshal(source, target)
	}
	return fmt.Errorf("Invalid body format: %s", format)
}

func IsEqual(rVal, eVal interface{}) bool {
	result := (rVal == eVal)
	if result {
		return result
	}
	rStr := fmt.Sprintf("%v", rVal)
	eStr := fmt.Sprintf("%v", eVal)
	return rStr == eStr
}

func VariableInfo(label string, val interface{}) {
	fmt.Printf(" - %s: [%v], type: %s\n", label, val, reflect.ValueOf(val).Type().String())
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
	IsEqualTo *int `yaml:"is-equal-to,omitempty" json:"is-equal-to"`
	BelongsTo []int `yaml:"belongs-to,omitempty" json:"belongs-to"`
}

type MeasureHeaders struct {
	Total *MeasureNumber `yaml:"total,omitempty" json:"total"`
	Items []MeasureHeader `yaml:"items,omitempty" json:"items"`
}

type MeasureNumber struct {
	IsEqualTo *int `yaml:"is-equal-to,omitempty" json:"is-equal-to"`
	IsLT *int `yaml:"is-lt,omitempty" json:"is-lt"`
	IsLTE *int `yaml:"is-lte,omitempty" json:"is-lte"`
	IsGT *int `yaml:"is-gt,omitempty" json:"is-gt"`
	IsGTE *int `yaml:"is-gte,omitempty" json:"is-gte"`
}

type MeasureHeader struct {
	Name *string `yaml:"name" json:"name"`
	IsEqualTo *string `yaml:"is-equal-to" json:"is-equal-to"`
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
	Value interface{} `yaml:"value,omitempty" json:"value"`
}

type ExaminationResult struct {
	Duration time.Duration
	Errors map[string]error
	Response *HttpResponse
	Status string
}

type DiffReporter struct {
	path  cmp.Path
	diffs []string
}

func (r *DiffReporter) PushStep(ps cmp.PathStep) {
	r.path = append(r.path, ps)
}

func (r *DiffReporter) Report(rs cmp.Result) {
	if !rs.Equal() {
		vx, vy := r.path.Last().Values()
		if !IsZero(vx) {
			r.diffs = append(r.diffs, fmt.Sprintf("%#v:\n\t-: %+v\n\t+: %+v\n", r.path, vx, vy))
		}
	}
}

func (r *DiffReporter) PopStep() {
	r.path = r.path[:len(r.path)-1]
}

func (r *DiffReporter) String() string {
	return strings.Join(r.diffs, "\n")
}

func (r *DiffReporter) HasDiffs() bool {
	return len(r.diffs) > 0
}

func (r *DiffReporter) GetDiffs() []string {
	return r.diffs
}

func IsZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
