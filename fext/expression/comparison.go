package expression

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/fextpkg/cli/fext/config"
)

// this functions implemented cause golang not support __eq__ methods as a python

// == (equal for strings)
func eqs(a, b string) bool { return a == b }

// != (not equal for strings)
func nes(a, b string) bool { return a != b }

// && (and)
func and(a, b bool) bool { return a && b }

// || (or)
func or(a, b bool) bool { return a || b }

func getStrCompareFunc(operator string) (func(a string, b string) bool, error) {
	switch operator {
	case "==":
		return eqs, nil
	case "!=":
		return nes, nil
	default:
		return nil, errors.New("Unsupported operator for strings: " + operator)
	}
}

// Returns indexes of the deepest and closet pair of brackets
func getBracketIndexes(s string) (int, int) {
	var start int
	for i, v := range s {
		if v == 40 { // 40 == (
			start = i
		} else if v == 41 { // 41 == )
			return start, i
		}
	}
	return -1, -1
}

func compareBool(a, b bool, operator string) (bool, error) {
	if operator == "and" {
		return and(a, b), nil
	} else if operator == "or" {
		return or(a, b), nil
	}

	return false, errors.New("Unknown operator :" + operator)
}

// Split comparison and logical operators
func splitExpOperators(exp string) ([]string, []string) {
	var comparison, logical []string
	re := regexp.MustCompile(`(\S+ [><=!]=? \S+|true|false)`)
	comparison = re.FindAllString(exp, -1)
	for _, v := range strings.Split(exp, " ") {
		if v == "and" || v == "or" {
			logical = append(logical, v)
		}
	}

	return comparison, logical
}

func compareSubExpression(s string) (bool, error) {
	defer func() { recover() }() // TODO signal about error
	if s == "" {
		return true, nil
	}
	var cResults []bool
	c, l := splitExpOperators(s)
	for _, v := range c {
		comp := strings.Split(v, " ") // [value, operator, value]
		if len(comp) > 1 {
			// remove quotes
			comp[2] = comp[2][1 : len(comp[2])-1]
		} else { // bool result from past operations
			value, err := strconv.ParseBool(comp[0])
			if err != nil {
				return false, err
			}
			cResults = append(cResults, value)
			continue
		}

		// set markers value (PEP 508)
		switch comp[0] {
		case "python_version", "python_full_version":
			value, err := CompareVersion(config.PythonVersion, comp[1], comp[2])
			if err != nil {
				return false, err
			}
			cResults = append(cResults, value)
		case "sys_platform":
			compareFunc, err := getStrCompareFunc(comp[1])
			if err != nil {
				return false, err
			}
			cResults = append(cResults, compareFunc(config.SysPlatform, comp[2]))
		case "extra":
			cResults = append(cResults, true)
		default: // skip unknown marker
			cResults = append(cResults, false)
		}
	}

	var lastResult bool
	lastResult = cResults[0]
	for i, v := range l {
		var err error
		lastResult, err = compareBool(lastResult, cResults[i+1], v)
		if err != nil {
			return false, err
		}
	}

	return lastResult, nil
}

func CompareExpression(exp string) (bool, error) {
	var s, e int // indexes
	var sub string
	for {
		s, e = getBracketIndexes(exp)
		if s != -1 {
			sub = exp[s : e+1]        // +1 for collect close bracket
			sub = sub[1 : len(sub)-1] // remove brackets
			r, err := compareSubExpression(sub)
			if err != nil {
				return false, err
			}
			exp = exp[:s] + strconv.FormatBool(r) + exp[e+1:]
		} else {
			return compareSubExpression(exp)
		}
	}
}
