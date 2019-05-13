package script

import (
	"regexp"
	"strings"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/utils"
)


type SelectorOptions interface {
	GetTestName() string
}

type Selector struct{
	testName string
	testNameRe *regexp.Regexp
}

func NewSelector(opts SelectorOptions) (ref *Selector, err error) {
	var testName string
	if opts != nil {
		testName = opts.GetTestName()
	}

	var testNameRe *regexp.Regexp
	if len(testName) > 0 {
		if utils.TEST_CASE_TITLE_REGEXP.MatchString(testName) {
			testName = standardizeName(testName)
		} else {
			re, err := regexp.Compile(strings.ToLower(testName))
			if err == nil {
				testNameRe = re
			}
		}
	}

	ref = &Selector{ testName: testName, testNameRe: testNameRe }

	return ref, err
}

func (r *Selector) TypeOfTestNameFilter() string {
	more := "blank"
	if len(r.testName) > 0 {
		more = "string"
	}
	if r.testNameRe != nil {
		more = "regexp"
	}
	return more
}

func (r *Selector) GetTestNameFilter() string {
	return r.testName
}

func (r *Selector) IsMatched(testName string) bool {
	if len(r.testName) == 0 {
		return true
	} else {
		name := standardizeName(testName)
		if r.testNameRe == nil {
			if strings.Contains(name, r.testName) {
				return true
			}
		} else {
			if r.testNameRe.MatchString(name) {
				return true
			}
		}
	}
	return false
}

func (r *Selector) GetTestCases(descriptors map[string]*Descriptor) []*engine.TestCase {
	testcases := make([]*engine.TestCase, 0)
	for _, d := range descriptors {
		testsuite := d.TestSuite
		if testsuite != nil {
			for _, testcase := range testsuite.TestCases {
				if testcase != nil && r.IsMatched(testcase.Title) {
					testcases = append(testcases, testcase)
				}
			}
		}
	}
	return testcases
}

func standardizeName(name string) string {
	if len(name) == 0 {
		return name
	}
	name = strings.ToLower(name)
	name = strings.Join(strings.Fields(strings.TrimSpace(name)), " ")
	return name
}
