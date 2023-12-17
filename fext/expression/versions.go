package expression

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/fextpkg/cli/fext/ferror"
)

var (
	operatorList = []rune{'>', '<', '=', '!'}
)

type Condition struct {
	Value    string
	Operator string
}

// parseVersion parses and splits the version into four semantic parts:
// major, minor, patch, and pre. It returns an array of the first three parts
// and separately returns the pre-version.
// If the pre-version is not provided, it will return 0. If any part is missing,
// it will be replaced with 0. Therefore, if an empty string is passed,
// an array of zeros will be returned. It throws an error if there is a syntax
// violation (error during the conversion to a number).
func parseVersion(s string) ([3]int, int, error) {
	var output [3]int
	var pre int

	parts := strings.Split(s, ".")
	partsCount := len(parts)

	for i := 0; i < 3; i++ {
		// Iterate only over existing elements and replace the rest with zeros
		if i < partsCount {
			value := parts[i]
			intValue, err := strconv.Atoi(value)
			if err != nil {
				// We receive an error if the string contains letters.
				// This means that the loop has reached the patch version and
				// it contains the pre-version
				intValue, pre, err = parsePreVersion(value)
				if err != nil {
					return output, 0, err
				}
			}
			output[i] = intValue
		} else {
			output[i] = 0
		}
	}

	return output, pre, nil
}

// parsePreVersion parses the version segment that contains letters and splits
// it into two parts: the patch version and the pre-version.
// It accepts the patch version as an argument and returns the split version numbers.
// If a regular number (without letters) or an empty string is passed, it will
// return two zeros without any errors.
func parsePreVersion(s string) (int, int, error) {
	var patchValue, preValue int
	var err error

	for i, v := range s {
		// First, we search for the index of the letter
		if !unicode.IsDigit(v) {
			// Next, we convert everything before it into a number
			patchValue, err = strconv.Atoi(s[:i])
			// TODO: replace this hack with a proper check for alpha, beta,
			// or release candidate
			if err != nil { // unknown error
				return 0, 0, err
			}
			// Then, we utilize a small hack and convert the entire string
			// after the letter into a number
			preValue = getStringIndexSum(s[i:])
			break
		}
	}

	return patchValue, preValue, nil
}

// getStringIndexSum converts each character from the string into its index
// and returns the sum of the indices.
func getStringIndexSum(s string) int {
	var output int
	for _, v := range s {
		output += int(v)
	}

	return output
}

// compareVersion returns the comparison result between versions.
// If a > b, it returns 1. If a < b, it returns -1. If a == b, it returns 0.
func compareVersion(a, b string) (int, error) {
	v1, v1pre, err := parseVersion(a)
	if err != nil {
		return 0, err
	}
	v2, v2pre, err := parseVersion(b)
	if err != nil {
		return 0, err
	}

	// First, we compare the first three semantic parts
	for i := 0; i < 3; i++ {
		if v1[i] > v2[i] {
			return 1, nil
		} else if v2[i] > v1[i] {
			return -1, nil
		}
	}

	// If the above comparison fails, compare the pre-version
	if v1pre&v2pre != 0 {
		if v1pre > v2pre {
			return 1, nil
		} else if v2pre > v1pre {
			return -1, nil
		}
	}

	// All values are identical to each other
	return 0, nil
}

// CompareVersion works by means of comparing each part of the version. The
// version is divided into semantic parts (major, minor, patch, pre) and
// converted into numbers that are compared one after another. For example,
// 4.0.0a >= 4.0.0rc2 will return the result false because alpha build has less
// weight than release candidate build.
//
// The version is parsed according to the PEP 440 standard, which means that the
// semantic version cannot be compared using it. Comparison of pre-release
// versions occurs only if they are specified in both versions. Otherwise, the
// check will be skipped
func CompareVersion(v1, op, v2 string) (bool, error) {
	res, err := compareVersion(v1, v2)
	if err != nil {
		return false, err
	}

	if res < 0 && (op == "<" || op == "<=" || op == "!=") {
		return true, nil
	} else if res > 0 && (op == ">" || op == ">=" || op == "!=") {
		return true, nil
	} else if res == 0 && (op == "==" || op == ">=" || op == "<=") {
		return true, nil
	} else if !strings.ContainsAny(op, "><=!") {
		return false, &ferror.UnexpectedOperator{Operator: op}
	}

	return false, nil
}

// isOperator checks if a rune corresponds to a specific operator symbol.
func isOperator(char rune) bool {
	for _, op := range operatorList {
		if op == char {
			return true
		}
	}

	return false
}

// splitConditions separates the comparison operator from the value and
// combines them. Returns a list of operators with their values.
func splitConditions(exp string) []Condition {
	var cond []Condition
	var op, version strings.Builder

	// Use a small trick by adding the operator to the end of the string.
	// This helps to avoid unnecessary code snippets.
	for _, char := range exp + "<" {
		if isOperator(char) {
			// Checking if the operator is the first one in the expression
			if version.Len() != 0 {
				cond = append(cond, Condition{
					Value:    version.String(),
					Operator: op.String(),
				})
				op.Reset()
				version.Reset()
			}

			op.WriteRune(char)
		} else {
			version.WriteRune(char)
		}
	}

	return cond
}

// ParseConditions separates the package name from the operators.
// The expressions must not contain spaces. Returns the package name and a list
// of operators with their values.
func ParseConditions(exp string) (string, []Condition) {
	for i, char := range exp {
		// Iterate through the string in search of an operator
		if isOperator(char) {
			// Split the package name and obtain the operators for the query
			return exp[:i], splitConditions(exp[i:])
		}
	}

	return exp, nil
}
