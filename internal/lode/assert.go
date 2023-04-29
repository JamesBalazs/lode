package lode

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"regexp"
	"strconv"
	"strings"
)

type Assertion interface {
	assert(lode *Lode) bool
}

func ParseAssertionArray(list []interface{}) []Assertion {
	var assertions []Assertion
	for _, m := range list {
		assertion := m.(map[string]interface{})
		assertions = append(assertions, NewAssertion(assertion))
	}
	return assertions
}

func NewAssertion(m map[string]interface{}) Assertion {
	typ, present := m["type"].(string)
	if !present {
		str, _ := json.Marshal(m)
		panic(fmt.Sprintf("no type provided in assertion: %s", str))
	}
	switch typ {
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
	case "greaterThan":
		return NewGreaterThanAssertion(m)
	case "lessThan":
		return NewLessThanAssertion(m)
	case "not":
		return NewNotAssertion(m)
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

func NewGreaterThanAssertion(m map[string]interface{}) GreaterThanAssertion {
	propertyName := m["property"].(string)
	subPropertyName, _ := m["key"].(string)
	greaterThan := cast.ToFloat64(m["greaterThan"]) // TODO switch to ToFloat64E and handle
	return GreaterThanAssertion{propertyName: propertyName, subPropertyName: subPropertyName, greaterThan: greaterThan}
}

type GreaterThanAssertion struct {
	propertyName    string
	subPropertyName string
	greaterThan     float64
}

func (a GreaterThanAssertion) assert(lode *Lode) bool {
	result := true

	switch a.propertyName {
	case "body":
		for _, response := range lode.ResponseTimings {
			if bodyF, err := strconv.ParseFloat(response.Response.Body, 64); err != nil {
				result = false
			} else {
				result = result && bodyF > a.greaterThan
			}
		}
	case "headers":
		if len(a.subPropertyName) == 0 {
			panic("no header key provided")
		}

		for _, response := range lode.ResponseTimings {
			header := response.Response.Header.HttpHeader.Get(a.subPropertyName)
			if headerF, err := strconv.ParseFloat(header, 64); err != nil {
				result = false
			} else {
				result = result && headerF > a.greaterThan
			}
		}
	case "status":
		for _, response := range lode.ResponseTimings {
			result = result && float64(response.Response.StatusCode) > a.greaterThan
		}
	default:
		panic("unknown property in assertion")
	}

	return result
}

func NewLessThanAssertion(m map[string]interface{}) LessThanAssertion {
	propertyName := m["property"].(string)
	subPropertyName, _ := m["key"].(string)
	lessThan := cast.ToFloat64(m["lessThan"]) // TODO switch to ToFloat64E and handle
	return LessThanAssertion{propertyName: propertyName, subPropertyName: subPropertyName, lessThan: lessThan}
}

type LessThanAssertion struct {
	propertyName    string
	subPropertyName string
	lessThan        float64
}

func (a LessThanAssertion) assert(lode *Lode) bool {
	result := true

	switch a.propertyName {
	case "body":
		for _, response := range lode.ResponseTimings {
			if bodyF, err := strconv.ParseFloat(response.Response.Body, 64); err != nil {
				result = false
			} else {
				result = result && bodyF < a.lessThan
			}
		}
	case "headers":
		if len(a.subPropertyName) == 0 {
			panic("no header key provided")
		}

		for _, response := range lode.ResponseTimings {
			header := response.Response.Header.HttpHeader.Get(a.subPropertyName)
			if headerF, err := strconv.ParseFloat(header, 64); err != nil {
				result = false
			} else {
				result = result && headerF < a.lessThan
			}
		}
	case "status":
		for _, response := range lode.ResponseTimings {
			result = result && float64(response.Response.StatusCode) < a.lessThan
		}
	default:
		panic("unknown property in assertion")
	}

	return result
}

func NewNotAssertion(m map[string]interface{}) NotAssertion {
	assertion := m["assertion"].(map[string]interface{})
	return NotAssertion{assertion: NewAssertion(assertion)}
}

type NotAssertion struct {
	assertion Assertion
}

func (a NotAssertion) assert(lode *Lode) bool {
	return !a.assertion.assert(lode)
}
