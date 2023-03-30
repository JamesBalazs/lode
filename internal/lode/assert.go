package lode

import (
	"fmt"
	"regexp"
	"strings"
)

type Assertion interface {
	assert(lode *Lode) bool
}

func ParseAssertionArray(list []interface{}) []Assertion {
	assertions := make([]Assertion, len(list))
	for _, m := range list {
		assertion := m.(map[string]interface{})
		assertions = append(assertions, NewAssertion(assertion))
	}
	return assertions
}

func NewAssertion(m map[string]interface{}) Assertion {
	switch m["type"].(string) {
	case "and":
		return NewAndAssertion(m)
	case "or":
		return NewOrAssertion(m)
	case "equals":
		return NewEqualsAssertion(m)
	case "contains":
		return NewContainsAssertion(m)
	case "matches":
		return NewMatchesAssertion(m)
	default:
		panic("unknown assertion type")
	}
}

func NewAndAssertion(m map[string]interface{}) AndAssertion {
	assertions := m["assertions"].([]interface{})
	list := ParseAssertionArray(assertions)
	return AndAssertion{list: list}
}

type AndAssertion struct {
	list []Assertion
}

func (a AndAssertion) assert(lode *Lode) bool {
	result := true
	for _, assertion := range a.list {
		result = result && assertion.assert(lode)
	}
	return result
}

func NewOrAssertion(m map[string]interface{}) OrAssertion {
	assertions := m["assertions"].([]interface{})
	list := ParseAssertionArray(assertions)
	return OrAssertion{list: list}
}

type OrAssertion struct {
	list []Assertion
}

func (a OrAssertion) assert(lode *Lode) bool {
	result := false
	for _, assertion := range a.list {
		result = result || assertion.assert(lode)
	}
	return result
}

func NewEqualsAssertion(m map[string]interface{}) EqualsAssertion {
	propertyName := m["property"].(string)
	subPropertyName, _ := m["key"].(string)
	equals := m["equals"].(string)
	return EqualsAssertion{propertyName: propertyName, subPropertyName: subPropertyName, equals: equals}
}

type EqualsAssertion struct {
	propertyName    string
	subPropertyName string
	equals          string
}

func (a EqualsAssertion) assert(lode *Lode) bool {
	result := true

	switch a.propertyName {
	case "body":
		for _, response := range lode.ResponseTimings {
			result = result && response.Response.Body == a.equals
		}
	case "headers":
		if len(a.subPropertyName) == 0 {
			panic("no header key provided")
		}

		for _, response := range lode.ResponseTimings {
			result = result && response.Response.Header.HttpHeader.Get(a.subPropertyName) == a.equals
		}
	case "status":
		for _, response := range lode.ResponseTimings {
			result = result && fmt.Sprint(response.Response.StatusCode) == a.equals
		}
	default:
		panic("unknown property in assertion")
	}

	return result
}

func NewContainsAssertion(m map[string]interface{}) ContainsAssertion {
	propertyName := m["property"].(string)
	subPropertyName, _ := m["key"].(string)
	contains := m["contains"].(string)
	return ContainsAssertion{propertyName: propertyName, subPropertyName: subPropertyName, contains: contains}
}

type ContainsAssertion struct {
	propertyName    string
	subPropertyName string
	contains        string
}

func (a ContainsAssertion) assert(lode *Lode) bool {
	result := true

	switch a.propertyName {
	case "body":
		for _, response := range lode.ResponseTimings {
			result = result && strings.Contains(response.Response.Body, a.contains)
		}
	case "headers":
		if len(a.subPropertyName) == 0 {
			panic("no header key provided")
		}

		for _, response := range lode.ResponseTimings {
			result = result && strings.Contains(response.Response.Header.HttpHeader.Get(a.subPropertyName), a.contains)
		}
	default:
		panic("unknown property in assertion")
	}

	return result
}

func NewMatchesAssertion(m map[string]interface{}) MatchesAssertion {
	propertyName := m["property"].(string)
	subPropertyName, _ := m["key"].(string)
	regexpString := m["regexp"].(string)
	if expression, err := regexp.Compile(regexpString); err != nil {
		panic(fmt.Sprintf("invalid matcher regex %s - error: %s", regexpString, err.Error()))
	} else {
		return MatchesAssertion{propertyName: propertyName, subPropertyName: subPropertyName, expression: expression}
	}
}

type MatchesAssertion struct {
	propertyName    string
	subPropertyName string
	expression      *regexp.Regexp
}

func (a MatchesAssertion) assert(lode *Lode) bool {
	result := true

	switch a.propertyName {
	case "body":
		for _, response := range lode.ResponseTimings {
			result = result && a.expression.MatchString(response.Response.Body)
		}
	case "headers":
		if len(a.subPropertyName) == 0 {
			panic("no header key provided")
		}

		for _, response := range lode.ResponseTimings {
			result = result && a.expression.MatchString(response.Response.Header.HttpHeader.Get(a.subPropertyName))
		}
	case "status":
		for _, response := range lode.ResponseTimings {
			result = result && a.expression.MatchString(fmt.Sprint(response.Response.StatusCode))
		}
	default:
		panic("unknown property in assertion")
	}

	return result
}

func a() {
}
