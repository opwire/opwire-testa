package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"gopkg.in/yaml.v2"
	"github.com/opwire/opwire-testa/lib/utils"
)

type TestGenerator struct {
	ExcludedHeaders []string
	Version string
}

func NewTestGenerator() (*TestGenerator, error) {
	ref := new(TestGenerator)
	ref.ExcludedHeaders = []string {
		"content-length",
		"date",
		"x-exec-duration",
	}
	return ref, nil
}

func (g *TestGenerator) generateTestCase(w io.Writer, req *HttpRequest, res *HttpResponse) error {
	s := TestCase{}
	s.Title = "<Generated testcase>"
	s.Version = utils.RefOfString(g.Version)
	s.Request = req
	s.Expectation = g.generateExpectation(res)

	r := &GeneratedSnapshot{}
	r.TestCases = []TestCase{s}
	script, err := yaml.Marshal(r)
	if err != nil {
		fmt.Fprintln(w, "Cannot marshal generated testcase, error: %s", err)
		return err
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, string(script))

	return nil
}

func (g *TestGenerator) generateExpectation(res *HttpResponse) *Expectation {
	if res == nil {
		return nil
	}
	e := &Expectation{}

	// status-code
	sc := res.StatusCode
	e.StatusCode = &MeasureStatusCode{
		IsEqualTo: &sc,
	}

	// header
	total := len(res.Header)
	if total > 0 {
		e.Headers = &MeasureHeaders{
			HasTotal: &total,
			Items: make([]MeasureHeader, 0),
		}
		for key, vals := range res.Header {
			if utils.ContainsInsensitiveCase(g.ExcludedHeaders, key) {
				continue
			}
			if len(vals) == 1 {
				name := key
				value := vals[0]
				one := MeasureHeader{
					Name: &name,
					IsEqualTo: &value,
				}
				e.Headers.Items = append(e.Headers.Items, one)
			}
		}
	}

	// body
	e.Body = &MeasureBody{}

	obj := make(map[string]interface{}, 0)
	if e.Body.HasFormat == nil {
		if err := json.Unmarshal(res.Body, &obj); err == nil {
			e.Body.HasFormat = utils.RefOfString("json")
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
			e.Body.HasFormat = utils.RefOfString("yaml")
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
		e.Body.HasFormat = utils.RefOfString("flat")
		e.Body.IsEqualTo = utils.RefOfString(string(res.Body))
		e.Body.MatchWith = utils.RefOfString(".*")
	}

	return e
}
