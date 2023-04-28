package lode

import (
	"encoding/json"
	"github.com/JamesBalazs/lode/internal/responseTimings"
	"github.com/stretchr/testify/assert"
	"net/http"
	"regexp"
	"testing"
)

type passingAssertion struct{}

func (passingAssertion) assert(lode *Lode) bool { return true }

type failingAssertion struct{}

func (failingAssertion) assert(lode *Lode) bool { return false }

const testAssertion = `{"type": "and", "assertions": [
	{"type": "or", "assertions": [
		{"type": "equals", "property": "body", "equals": "this"},
		{"type": "contains", "property": "body", "contains": "that"}
	]},
	{"type": "matches", "property": "status", "regexp": "2\\d\\d"},
	{"type": "greaterThan", "property": "headers", "key": "X-Custom", "greaterThan": 199},
	{"type": "lessThan", "property": "headers", "key": "X-Custom", "lessThan": "500"}
]}`

func TestNewAssertion(t *testing.T) {
	assertionMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(testAssertion), &assertionMap)
	assert.Nil(t, err)

	assertion := NewAssertion(assertionMap)

	assert.Equal(t, AndAssertion{
		list: []Assertion{
			OrAssertion{list: []Assertion{
				EqualsAssertion{propertyName: "body", equals: "this"},
				ContainsAssertion{propertyName: "body", contains: "that"},
			}},
			MatchesAssertion{propertyName: "status", expression: regexp.MustCompile("2\\d\\d")},
			GreaterThanAssertion{propertyName: "headers", subPropertyName: "X-Custom", greaterThan: 199},
			LessThanAssertion{propertyName: "headers", subPropertyName: "X-Custom", lessThan: 500},
		},
	}, assertion)
}

func TestNewAssertion_NoType(t *testing.T) {
	assertionMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(`{"something": "else"}`), &assertionMap)
	assert.Nil(t, err)

	assert.Panics(t, func() {
		NewAssertion(assertionMap)
	})
}

func TestNewAssertion_UnknownType(t *testing.T) {
	assertionMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(`{"type": "unknown"}`), &assertionMap)
	assert.Nil(t, err)

	assert.Panics(t, func() {
		NewAssertion(assertionMap)
	})
}

func TestParseAssertionArray(t *testing.T) {
	arr := []interface{}{
		map[string]interface{}{
			"type":     "equals",
			"property": "body",
			"equals":   "this",
		},
		map[string]interface{}{
			"type":     "equals",
			"property": "body",
			"equals":   "that",
		},
	}

	result := ParseAssertionArray(arr)

	assert.Equal(t, []Assertion{
		EqualsAssertion{propertyName: "body", subPropertyName: "", equals: "this"},
		EqualsAssertion{propertyName: "body", subPropertyName: "", equals: "that"},
	}, result)
}

func TestNewAndAssertion(t *testing.T) {
	assertion := map[string]interface{}{
		"type": "and",
		"assertions": []interface{}{
			map[string]interface{}{
				"type":     "equals",
				"property": "body",
				"equals":   "this",
			},
			map[string]interface{}{
				"type":     "equals",
				"property": "body",
				"equals":   "that",
			},
		},
	}

	result := NewAndAssertion(assertion)

	assert.Equal(t, AndAssertion{
		list: []Assertion{
			EqualsAssertion{
				propertyName: "body",
				equals:       "this",
			},
			EqualsAssertion{
				propertyName: "body",
				equals:       "that",
			},
		},
	}, result)
}

func TestAndAssertion_assert(t *testing.T) {
	lode := &Lode{}
	assertion := AndAssertion{
		list: []Assertion{
			failingAssertion{},
			failingAssertion{},
		},
	}

	assert.False(t, assertion.assert(lode))

	assertion = AndAssertion{
		list: []Assertion{
			passingAssertion{},
			failingAssertion{},
		},
	}

	assert.False(t, assertion.assert(lode))

	assertion = AndAssertion{
		list: []Assertion{
			failingAssertion{},
			passingAssertion{},
		},
	}

	assert.False(t, assertion.assert(lode))

	assertion = AndAssertion{
		list: []Assertion{
			passingAssertion{},
			passingAssertion{},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestNewOrAssertion(t *testing.T) {
	assertion := map[string]interface{}{
		"type": "or",
		"assertions": []interface{}{
			map[string]interface{}{
				"type":     "equals",
				"property": "body",
				"equals":   "this",
			},
			map[string]interface{}{
				"type":     "equals",
				"property": "body",
				"equals":   "that",
			},
		},
	}

	result := NewOrAssertion(assertion)

	assert.Equal(t, OrAssertion{
		list: []Assertion{
			EqualsAssertion{
				propertyName: "body",
				equals:       "this",
			},
			EqualsAssertion{
				propertyName: "body",
				equals:       "that",
			},
		},
	}, result)
}

func TestOrAssertion_assert(t *testing.T) {
	lode := &Lode{}
	assertion := OrAssertion{
		list: []Assertion{
			failingAssertion{},
			failingAssertion{},
		},
	}

	assert.False(t, assertion.assert(lode))

	assertion = OrAssertion{
		list: []Assertion{
			failingAssertion{},
			passingAssertion{},
		},
	}

	assert.True(t, assertion.assert(lode))

	assertion = OrAssertion{
		list: []Assertion{
			passingAssertion{},
			failingAssertion{},
		},
	}

	assert.True(t, assertion.assert(lode))

	assertion = OrAssertion{
		list: []Assertion{
			passingAssertion{},
			passingAssertion{},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestNewEqualsAssertion(t *testing.T) {
	assertion := map[string]interface{}{
		"type":     "equals",
		"property": "headers",
		"key":      "Header-Name",
		"equals":   "this",
	}

	result := NewEqualsAssertion(assertion)

	assert.Equal(t, EqualsAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		equals:          "this",
	}, result)
}

func TestEqualsAssertion_assertBody(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "this"}},
			{Response: &responseTimings.Response{Body: "that"}},
		},
	}

	assertion := EqualsAssertion{
		propertyName: "body",
		equals:       "this",
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "this"}},
			{Response: &responseTimings.Response{Body: "this"}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestEqualsAssertion_assertHeaders(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this"}}}}},
		},
	}

	assertion := EqualsAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		equals:          "this",
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this"}}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this"}}}}},
		},
	}

	assert.True(t, assertion.assert(lode))

	assertion = EqualsAssertion{
		propertyName: "headers",
		equals:       "this",
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestEqualsAssertion_assertStatus(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 200}},
			{Response: &responseTimings.Response{StatusCode: 201}},
		},
	}

	assertion := EqualsAssertion{
		propertyName: "status",
		equals:       "201",
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 201}},
			{Response: &responseTimings.Response{StatusCode: 201}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestEqualsAssertion_assertUnknownProperty(t *testing.T) {
	lode := &Lode{}
	assertion := EqualsAssertion{
		propertyName: "unknown",
		equals:       "201",
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestNewContainsAssertion(t *testing.T) {
	assertion := map[string]interface{}{
		"type":     "equals",
		"property": "headers",
		"key":      "Header-Name",
		"contains": "this",
	}

	result := NewContainsAssertion(assertion)

	assert.Equal(t, ContainsAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		contains:        "this",
	}, result)
}

func TestContainsAssertion_assertBody(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "this that"}},
			{Response: &responseTimings.Response{Body: "that"}},
		},
	}

	assertion := ContainsAssertion{
		propertyName: "body",
		contains:     "this",
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "this that"}},
			{Response: &responseTimings.Response{Body: "that this"}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestContainsAssertion_assertHeaders(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this that"}}}}},
		},
	}

	assertion := ContainsAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		contains:        "this",
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this that"}}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"that this"}}}}},
		},
	}

	assert.True(t, assertion.assert(lode))

	assertion = ContainsAssertion{
		propertyName: "headers",
		contains:     "this",
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestContainsAssertion_assertUnknownProperty(t *testing.T) {
	lode := &Lode{}
	assertion := ContainsAssertion{
		propertyName: "unknown",
		contains:     "this",
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestNewMatchesAssertion(t *testing.T) {
	assertion := map[string]interface{}{
		"type":     "matches",
		"property": "headers",
		"key":      "Header-Name",
		"regexp":   "some.*",
	}

	result := NewMatchesAssertion(assertion)

	assert.Equal(t, MatchesAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		expression:      regexp.MustCompile("some.*"),
	}, result)

	assertion = map[string]interface{}{
		"type":     "matches",
		"property": "headers",
		"key":      "Header-Name",
		"regexp":   "invalid\\",
	}

	assert.Panics(t, func() {
		NewMatchesAssertion(assertion)
	})
}

func TestMatchesAssertion_assertBody(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "this thing"}},
			{Response: &responseTimings.Response{Body: "that thing"}},
		},
	}

	assertion := MatchesAssertion{
		propertyName: "body",
		expression:   regexp.MustCompile("this.*"),
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "this thing"}},
			{Response: &responseTimings.Response{Body: "this other"}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestMatchesAssertion_assertHeaders(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this"}}}}},
		},
	}

	assertion := MatchesAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		expression:      regexp.MustCompile("this.*"),
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this thing"}}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"this other"}}}}},
		},
	}

	assert.True(t, assertion.assert(lode))

	assertion = MatchesAssertion{
		propertyName: "headers",
		expression:   regexp.MustCompile("this.*"),
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestMatchesAssertion_assertStatus(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 200}},
			{Response: &responseTimings.Response{StatusCode: 401}},
		},
	}

	assertion := MatchesAssertion{
		propertyName: "status",
		expression:   regexp.MustCompile("2\\d\\d"),
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 201}},
			{Response: &responseTimings.Response{StatusCode: 200}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestMatchesAssertion_assertUnknownProperty(t *testing.T) {
	lode := &Lode{}
	assertion := MatchesAssertion{
		propertyName: "unknown",
		expression:   regexp.MustCompile("2\\d\\d"),
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestNewGreaterThanAssertion(t *testing.T) {
	assertion := map[string]interface{}{
		"type":        "greaterThan",
		"property":    "headers",
		"key":         "Header-Name",
		"greaterThan": 200,
	}

	result := NewGreaterThanAssertion(assertion)

	assert.Equal(t, GreaterThanAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		greaterThan:     200,
	}, result)
}

func TestGreaterThanAssertion_assertBody(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "200"}},
			{Response: &responseTimings.Response{Body: "201"}},
		},
	}

	assertion := GreaterThanAssertion{
		propertyName: "body",
		greaterThan:  200,
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "abc"}},
			{Response: &responseTimings.Response{Body: "201"}},
		},
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "201"}},
			{Response: &responseTimings.Response{Body: "201"}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestGreaterThanAssertion_assertHeaders(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"200"}}}}},
		},
	}

	assertion := GreaterThanAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		greaterThan:     200,
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"201"}}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"201"}}}}},
		},
	}

	assert.True(t, assertion.assert(lode))

	assertion = GreaterThanAssertion{
		propertyName: "headers",
		greaterThan:  201,
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestGreaterThanAssertion_assertStatus(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 200}},
			{Response: &responseTimings.Response{StatusCode: 201}},
		},
	}

	assertion := GreaterThanAssertion{
		propertyName: "status",
		greaterThan:  200,
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 201}},
			{Response: &responseTimings.Response{StatusCode: 201}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestGreaterThanAssertion_assertUnknownProperty(t *testing.T) {
	lode := &Lode{}
	assertion := GreaterThanAssertion{
		propertyName: "unknown",
		greaterThan:  200,
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestNewLessThanAssertion(t *testing.T) {
	assertion := map[string]interface{}{
		"type":     "greaterThan",
		"property": "headers",
		"key":      "Header-Name",
		"lessThan": 300,
	}

	result := NewLessThanAssertion(assertion)

	assert.Equal(t, LessThanAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		lessThan:        300,
	}, result)
}

func TestLessThanAssertion_assertBody(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "301"}},
			{Response: &responseTimings.Response{Body: "204"}},
		},
	}

	assertion := LessThanAssertion{
		propertyName: "body",
		lessThan:     300,
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "abc"}},
			{Response: &responseTimings.Response{Body: "204"}},
		},
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Body: "299"}},
			{Response: &responseTimings.Response{Body: "200"}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestLessThanAssertion_assertHeaders(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"123"}}}}},
		},
	}

	assertion := LessThanAssertion{
		propertyName:    "headers",
		subPropertyName: "Header-Name",
		lessThan:        124,
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"123"}}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"abc"}}}}},
		},
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"123"}}}}},
			{Response: &responseTimings.Response{Header: responseTimings.Header{HttpHeader: http.Header{"Header-Name": []string{"101"}}}}},
		},
	}

	assert.True(t, assertion.assert(lode))

	assertion = LessThanAssertion{
		propertyName: "headers",
		lessThan:     123,
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}

func TestLessThanAssertion_assertStatus(t *testing.T) {
	lode := &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 200}},
			{Response: &responseTimings.Response{StatusCode: 201}},
		},
	}

	assertion := LessThanAssertion{
		propertyName: "status",
		lessThan:     201,
	}

	assert.False(t, assertion.assert(lode))

	lode = &Lode{
		ResponseTimings: []responseTimings.ResponseTiming{
			{Response: &responseTimings.Response{StatusCode: 200}},
			{Response: &responseTimings.Response{StatusCode: 100}},
		},
	}

	assert.True(t, assertion.assert(lode))
}

func TestLessThanAssertion_assertUnknownProperty(t *testing.T) {
	lode := &Lode{}
	assertion := LessThanAssertion{
		propertyName: "unknown",
		lessThan:     123,
	}

	assert.Panics(t, func() {
		assertion.assert(lode)
	})
}
