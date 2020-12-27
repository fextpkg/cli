package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// this functions implemented cause golang not support __eq__ methods as a python

// == (equal)
func eq(a, b int) bool { return a == b }
// == (equal for strings)
func eqs(a, b string) bool { return a == b }
// != (not equal)
func ne(a, b int) bool { return a != b }
// != (not equal for strings)
func nes(a, b string) bool { return a != b }
// <= (less then or equal)
func lte(a, b int) bool { return a <= b }
// < (less then)
func lt(a, b int) bool { return a < b }
// >= (greater then or equal)
func gte(a, b int) bool { return a >= b }
// > (greater then)
func gt(a, b int) bool { return a > b }
// && (and)
func and(a, b bool) bool { return a && b }
// || (or)
func or(a, b bool) bool { return a || b }

func getCompareFunc(operator string) (func(a int, b int) bool, error) {
	switch operator {
	case "==": return eq, nil
	case "!=": return ne, nil
	case "<=": return lte, nil
	case "<": return lt, nil
	case ">=": return gte, nil
	case ">": return gt, nil
	default:
		return nil, errors.New("Unknown operator: " + operator)
	}
}

func getStrCompareFunc(operator string) (func(a string, b string) bool, error) {
	switch operator {
	case "==": return eqs, nil
	case "!=": return nes, nil
	default:
		return nil, errors.New("Unsupported operator for strings: " + operator)
	}
}

func convertLettersToIntString(letters string) string {
	var out string

	for _, v := range letters {
		if !unicode.IsDigit(v) {
			// 97 is index of start letters in lowercase. We don't use original
			// indexes cause compare function works faster
			out += strconv.Itoa(int(v - 96))
		} else {
			out += string(v)
		}
	}

	return out
}

func compare(a, b int, operator string) (bool, error) {
	compareFunc, err := getCompareFunc(operator)
	if err != nil {
		return false, err
	}

	return compareFunc(a, b), nil
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

func Convert(s string) string {
	return convertLettersToIntString(s)
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
	re, _ := regexp.Compile(`(\S+ [><=!]=? \S+|true|false)`)
	comparison = re.FindAllString(exp, -1)
	for _, v := range strings.Split(exp, " ") {
		if v == "and" || v == "or" {
			logical = append(logical, v)
		}
	}

	return comparison, logical
}

// libDir used for set correct markers (if any)
func compareSubExpression(s, libDir string) (bool, error) {
	defer func() {recover()}() // TODO signal about error
	var cResults []bool
	c, l := splitExpOperators(s)
	for _, v := range c {
		comp := strings.Split(v, " ") // [value, operator, value]
		if len(comp) > 1 {
			// remove quotes
			comp[2] = comp[2][1:len(comp[2])-1]
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
		case "python_version":
			version := GetPythonVersion(libDir)
			value, err := CompareVersion(version, comp[1], comp[2])
			if err != nil {
				return false, err
			}
			cResults = append(cResults, value)
		case "sys_platform":
			platform := GetSysPlatform()
			compareFunc, err := getStrCompareFunc(comp[1])
			if err != nil {
				return false, err
			}
			cResults = append(cResults, compareFunc(platform, comp[2]))
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

func CompareExpression(exp, libDir string) (bool, error) {
	var s, e int // indexes
	var sub string
	for {
		s, e = getBracketIndexes(exp)
		if s != -1 {
			sub = exp[s:e + 1] // +1 for collect close bracket
			sub = sub[1:len(sub) - 1] // remove brackets
			r, err := compareSubExpression(sub, libDir)
			if err != nil {
				return false, err
			}
			exp = exp[:s] + strconv.FormatBool(r) + exp[e + 1:]
		} else {
			return compareSubExpression(exp, libDir)
		}
	}

	return false, nil
}

