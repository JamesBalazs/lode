package lode

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
	{"type": "matches", "property": "status", "regexp": "2\\d\\d"}
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
