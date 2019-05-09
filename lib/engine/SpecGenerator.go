package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
	"gopkg.in/yaml.v2"
	"github.com/opwire/opwire-testa/lib/utils"
)

type SpecGenerator struct {
	ExcludedHeaders []string
	Version string
}

func NewSpecGenerator() (*SpecGenerator, error) {
	ref := new(SpecGenerator)
	ref.ExcludedHeaders = []string {
		"content-length",
		"date",
		"x-exec-duration",
	}
	return ref, nil
}

func (g *SpecGenerator) generateTestCase(w io.Writer, req *HttpRequest, res *HttpResponse) error {
	s := TestCase{}
	s.Title = "<Generated testcase>"
	s.Version = utils.RefOfString(g.Version)
	s.Request = req
	s.Expectation = g.generateExpectation(res)
	s.CreatedTime = utils.RefOfString(time.Now().Format(time.RFC3339))
	s.Tags = []string {"snapshot"}
	username, err := utils.FindUsername()
	if err == nil {
		if tag, err := utils.StandardizeTagLabel(username); err == nil {
			s.Tags = append(s.Tags, tag)
		}
	}

	r := &GeneratedSnapshot{}
	r.TestCases = []TestCase{s}
	script, err := yaml.Marshal(r)
	if err != nil {
		fmt.Fprintf(w, "Cannot marshal generated testcase, error: %s\n", err)
		return err
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, string(script))

	return nil
}

func (g *SpecGenerator) generateExpectation(res *HttpResponse) *Expectation {
	if res == nil {
		return nil
	}
	e := &Expectation{}

	// status-code
	sc := res.StatusCode
	e.StatusCode = &MeasureStatusCode{
		Is: &ComparisonOperators{
			EqualTo: &sc,
			MemberOf: []interface{}{sc},
		},
	}

	// header
	total := len(res.Header)
	if total > 0 {
		e.Headers = &MeasureHeaders{
			Total: &MeasureTotal{
				Is: &ComparisonOperators{
					EqualTo: &total,
				},
			},
			Items: make([]MeasureHeader, 0),
		}
		count := 0
		for key, vals := range res.Header {
			if utils.ContainsInsensitiveCase(g.ExcludedHeaders, key) {
				continue
			}
			if len(vals) == 1 {
				name := key
				value := vals[0]
				one := MeasureHeader{
					Name: &name,
					Is: &ComparisonOperators{
						EqualTo: &value,
					},
				}
				e.Headers.Items = append(e.Headers.Items, one)
				count = count + 1
			}
		}
	}

	// body
	e.Body = &MeasureBody{}

	obj := make(map[string]interface{}, 0)
	if e.Body.HasFormat == nil {
		if err := json.Unmarshal(res.Body, &obj); err == nil {
			e.Body.HasFormat = utils.RefOfString(utils.BODY_FORMAT_JSON)
			var content string
			if out, err := json.MarshalIndent(obj, "", "  "); err == nil {
				content = string(out)
			} else {
				content = string(res.Body)
			}
			e.Body.Includes = &content
		}
	}

	if e.Body.HasFormat == nil {
		if err := yaml.Unmarshal(res.Body, &obj); err == nil {
			e.Body.HasFormat = utils.RefOfString(utils.BODY_FORMAT_YAML)
			var content string
			if out, err := yaml.Marshal(obj); err == nil {
				content = string(out)
			} else {
				content = string(res.Body)
			}
			e.Body.Includes = &content
		}
	}

	if e.Body.HasFormat == nil {
		e.Body.HasFormat = utils.RefOfString(utils.BODY_FORMAT_FLAT)
		e.Body.IsEqualTo = utils.RefOfString(string(res.Body))
		e.Body.MatchWith = utils.RefOfString(".*")
	}

	// body fields
	if len(obj) > 0 {
		flatten, _ := utils.Flatten("", obj)
		fields := make([]MeasureBodyField, 0)
		for key, val := range flatten {
			if val != nil {
				fields = append(fields, MeasureBodyField{
					Path: utils.RefOfString(key),
					Is: &ComparisonOperators{
						EqualTo: val,
					},
				})
			}
		}
		e.Body.Fields = fields
	}

	return e
}

type GeneratedSnapshot struct {
	TestCases []TestCase `yaml:"testcase-snapshot"`
}
