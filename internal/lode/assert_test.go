package lode

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testAssertion = `{"type": "and", "assertions": [
	{"type": "or", "assertions": [
		{"type": "equals", "property": "body", "equals": "this"},
		{"type": "equals", "property": "status", "equals": "201"}
	]},
	{"type": "contains", "property": "body", "contains": "that"}
]}`

func TestNewAssertion(t *testing.T) {
	assertionMap := map[string]interface{}{}
	json.Unmarshal([]byte(testAssertion), &assertionMap)

	assertion := NewAssertion(assertionMap)

	assert.Equal(t, AndAssertion{
		list: []Assertion{
			Assertion(nil),
			Assertion(nil),
			OrAssertion{list: []Assertion{
				Assertion(nil),
				Assertion(nil),
				EqualsAssertion{propertyName: "body", equals: "this"},
				EqualsAssertion{propertyName: "status", equals: "201"},
			}},
			ContainsAssertion{propertyName: "body", contains: "that"},
		},
	}, assertion)
}
