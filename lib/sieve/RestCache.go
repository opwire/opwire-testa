package sieve

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"github.com/opwire/opwire-testa/lib/client"
	"github.com/opwire/opwire-testa/lib/utils"
)

func NewRestCache() (*RestCache, error) {
	s := &RestCache{}
	s.restResult = make(map[string]*RestResult, 0)
	return s, nil
}

type RestCache struct {
	restResult map[string]*RestResult
}

func (s *RestCache) Query(query string) (string, error) {
	if s.restResult == nil {
		return utils.BLANK, fmt.Errorf("RestCache must be initialized")
	}

	q, err := Parse(query)
	if err != nil {
		return utils.BLANK, err
	}
	if q == nil {
		return utils.BLANK, fmt.Errorf("Query[%s] not found", query)
	}

	if len(q.TestID) == 0 {
		return utils.BLANK, fmt.Errorf("TestID must not be empty")
	}

	rr, ok := s.restResult[q.TestID]
	if !ok || rr == nil {
		return utils.BLANK, fmt.Errorf("RestResult[%s] not found", q.TestID)
	}

	switch(q.Attr) {
	case RESP_STATUS:
		return rr.Status, nil

	case RESP_STATUS_CODE:
		return fmt.Sprintf("%d", rr.StatusCode), nil

	case RESP_HEADER:
		if len(q.ItemKey) == 0 {
			return utils.BLANK, fmt.Errorf("Resp[%s].Header's name must be provided", q.TestID)
		}
		vals, found := rr.Header[q.ItemKey]
		if !found {
			if len(q.Default) > 0 {
				return q.Default, nil
			}
			return utils.BLANK, fmt.Errorf("Resp[%s].Header[%s] not found", q.TestID, q.ItemKey)
		}
		valsize := len(vals)
		if valsize == 0 {
			return utils.BLANK, fmt.Errorf("Resp[%s].Header[%s] is invalid", q.TestID, q.ItemKey)
		}
		if valsize == 1 {
			return vals[0], nil
		} else {
			return strings.Join(vals, ", "), nil
		}

	case RESP_BODY:
		if rr.Body == nil {
			if len(q.Default) > 0 {
				return q.Default, nil
			}
			return utils.BLANK, fmt.Errorf("Resp[%s].Body is nil", q.TestID)
		}
		return string(rr.Body), nil

	case RESP_BODY_FIELD:
		if rr.BodyField == nil {
			return utils.BLANK, fmt.Errorf("Resp[%s].BodyField must be initialized", q.TestID)
		}
		if len(q.ItemKey) == 0 {
			return utils.BLANK, fmt.Errorf("Resp[%s].BodyField's name must be provided", q.TestID)
		}
		val, found := rr.BodyField[q.ItemKey]
		if !found {
			if len(q.Default) > 0 {
				return q.Default, nil
			}
			return utils.BLANK, fmt.Errorf("Resp[%s].BodyField[%s] not found", q.TestID, q.ItemKey)
		}
		return fmt.Sprintf("%v", val), nil
	}
	return utils.BLANK, nil
}

func (s *RestCache) Get(testId string) (*RestResult, error) {
	if rr, ok := s.restResult[testId]; ok {
		return rr, nil
	} else {
		return nil, fmt.Errorf("RestResult[%s] not found", testId)
	}
}

func (s *RestCache) Store(testId string, res *client.HttpResponse) (*RestResult, error) {
	if s.restResult == nil {
		s.restResult = make(map[string]*RestResult, 0)
	}

	newRR, err := NewRestResult(res)
	if err != nil {
		return nil, err
	}

	oldRR, ok := s.restResult[testId]
	s.restResult[testId] = newRR

	if ok {
		return oldRR, nil
	}
	return nil, nil
}

func NewRestResult(lowRes *client.HttpResponse) (*RestResult, error) {
	if lowRes == nil {
		return nil, fmt.Errorf("HttpResponse must not be nil")
	}

	res := &RestResult{}

	res.Status = lowRes.Status
	res.StatusCode = lowRes.StatusCode
	res.Header = lowRes.Header
	res.ContentLength = lowRes.ContentLength
	res.Body = lowRes.Body

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
		res.BodyField = make(map[string]interface{})
		flatten, _ := utils.Flatten("", obj)
		for key, val := range flatten {
			res.BodyField[key] = val
		}
	}

	return res, nil
}

type RestResult struct {
	Status string
	StatusCode int
	Header http.Header
	ContentLength int64
	Body []byte
	BodyField map[string]interface{}
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
