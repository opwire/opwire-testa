package sieve

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"github.com/opwire/opwire-testa/lib/engine"
	"github.com/opwire/opwire-testa/lib/utils"
)

type RestCache struct {}

func NewRestCache() (*RestCache, error) {
	s := &RestCache{}
	return s, nil
}

func (s *RestCache) Read(query string) (string, error) {
	q, err := Parse(query)
	if err != nil {
		return "", err
	}

	_ = q
	return "", nil
}

func (s *RestCache) Save(query string) (string, error) {
	return "", nil
}

func NewRestResult(lowRes *engine.HttpResponse) (*RestResult, error) {
	if lowRes == nil {
		return nil, fmt.Errorf("HttpResponse must not be nil")
	}

	res := &RestResult{}

	res.Status = lowRes.Status
	res.StatusCode = lowRes.StatusCode
	res.Header = lowRes.Header
	res.ContentLength = lowRes.ContentLength
	res.Body = string(lowRes.Body)

	// BodyField
	obj := make(map[string]interface{}, 0)
	found := false

	// detect whether body has JSON format?
	if !found {
		if err := utils.Unmarshal(utils.BODY_FORMAT_JSON, lowRes.Body, &obj); err == nil {
			found = true
		}
	}

	if !found {
		if err := utils.Unmarshal(utils.BODY_FORMAT_YAML, lowRes.Body, &obj); err == nil {
			found = true
		}
	}

	if found {
		res.BodyField = make(map[string]string)
		flatten, _ := utils.Flatten("", obj)
		for key, val := range flatten {
			if val != nil {
				res.BodyField[key] = fmt.Sprintf("%v", val)
			}
		}
	}

	return res, nil
}

type RestResult struct {
	Status string
	StatusCode int
	Header http.Header
	ContentLength int64
	Body string
	BodyField map[string]string
}

type DataType int

const (
	_ DataType = iota
	RESP_STATUS
	RESP_STATUS_CODE
	RESP_HEADER
	RESP_BODY
	RESP_BODY_FIELD
)

type Query struct {
	TestID string
	Attr DataType
	ItemKey string
	Default string
}

var SIEVE_TESTCASE_STATUS_REGEXP = regexp.MustCompile(`(?i)^\s*case\[([^\]]*)\]\.Status\s*(\:\-([^\}]*))?\s*$`)
var SIEVE_TESTCASE_STATUS_CODE_REGEXP = regexp.MustCompile(`(?i)^\s*case\[([^\]]*)\]\.StatusCode\s*(\:\-([^\}]*))?\s*$`)
var SIEVE_TESTCASE_HEADER_REGEXP = regexp.MustCompile(`(?i)^\s*case\[([^\]]*)\]\.Header\[([^\]]*)\]\s*(\:\-([^\}]*))?\s*$`)
var SIEVE_TESTCASE_BODY_REGEXP = regexp.MustCompile(`(?i)^\s*case\[([^\]]*)\]\.Body\s*(\:\-([^\}]*))?\s*$`)
var SIEVE_TESTCASE_BODY_FIELD_REGEXP = regexp.MustCompile(`(?i)^\s*case\[([^\]]*)\]\.Body\[([^\]]*)\]\s*(\:\-([^\}]*))?\s*$`)

func Parse(query string) (*Query, error) {
	var q *Query
	q = extract2(RESP_STATUS, SIEVE_TESTCASE_STATUS_REGEXP.FindAllStringSubmatch(query, -1))
	if q != nil {
		return q, nil
	}
	q = extract2(RESP_STATUS_CODE, SIEVE_TESTCASE_STATUS_CODE_REGEXP.FindAllStringSubmatch(query, -1))
	if q != nil {
		return q, nil
	}
	q = extract3(RESP_HEADER, SIEVE_TESTCASE_HEADER_REGEXP.FindAllStringSubmatch(query, -1))
	if q != nil {
		return q, nil
	}
	q = extract2(RESP_BODY, SIEVE_TESTCASE_BODY_REGEXP.FindAllStringSubmatch(query, -1))
	if q != nil {
		return q, nil
	}
	q = extract3(RESP_BODY_FIELD, SIEVE_TESTCASE_BODY_FIELD_REGEXP.FindAllStringSubmatch(query, -1))
	if q != nil {
		return q, nil
	}
	return nil, nil
}

func extract2(attrType DataType, r [][]string) *Query {
	if r != nil && len(r) >= 1 && len(r[0]) >= 2 {
		ls := r[0]
		q := &Query{
			TestID: ls[1],
			Attr: attrType,
		}
		if len(ls) == 4 {
			q.Default = strings.TrimSpace(ls[3])
		}
		return q
	}
	return nil
}

func extract3(attrType DataType, r [][]string) *Query {
	if r != nil && len(r) >= 1 && len(r[0]) >= 3 {
		ls := r[0]
		q := &Query{
			TestID: ls[1],
			Attr: attrType,
			ItemKey: ls[2],
		}
		if len(ls) == 5 {
			q.Default = strings.TrimSpace(ls[4])
		}
		return q
	}
	return nil
}
