package expression

import (
	"strconv"
	"strings"

	"github.com/fextpkg/cli/fext/ferror"
)

type expression struct {
	v1 string
	v2 string
	op string
}

// CompareString compares the given values using the specified comparison
// operator for secure evaluation of Python string expressions.
// In case of an unexpected comparison operator, it returns an error.
func CompareString(a, operator, b string) (bool, error) {
	if operator == "==" {
		return a == b, nil
	} else if operator == "!=" {
		return a != b, nil
	} else {
		return false, &ferror.UnexpectedOperator{Operator: operator}
	}
}

// compareBool compares two boolean values using the specified operator.
// It is commonly used to evaluate and compare results from other expressions.
// It returns an error if an unexpected logical operator is passed.
func compareBool(a, b bool, operator string) (bool, error) {
	if operator == "and" {
		return a && b, nil
	} else if operator == "or" {
		return a || b, nil
	} else {
		return false, &ferror.UnexpectedOperator{Operator: operator}
	}
}

// Find the first deepest occurrence pair of parentheses.
// Return the indices or two -1 values if nothing is found.
func getBracketIndexes(s string) (int, int) {
	var start int
	for i, v := range s {
		if v == '(' {
			start = i
		} else if v == ')' {
			return start, i
		}
	}

	return -1, -1
}

// Parses expressions with a comparison operator and completed comparisons.
// Note that if the expression has already been compared (the string contains
// "true" or "false"), the attributes "v2" and "op" will be empty.
// Returns an error in case of syntax violation.
func parseExpressionWithOperators(s string) ([]expression, error) {
	var output []expression
	var exp expression

	delimitedString := strings.Split(s, " ")
	// Subtract one to enable comparison with the index
	length := len(delimitedString) - 1

	for i, sequence := range delimitedString {
		if strings.ContainsAny(sequence, "=><") {
			// Verify that there are elements on both sides
			if i >= length || i == 0 {
				// TODO: more readable
				return nil, ferror.SyntaxError
			}

			// The idea of the search is to find a comparison operator and
			// extract the values on its left and right sides
			exp = expression{
				v1: delimitedString[i-1],
				v2: delimitedString[i+1],
				op: sequence,
			}
			// TODO: here it may be worth to adding a check for the logical
			// operators or comparison operators on the left and right side.
			// In other words, if they are present, immediately return a SyntaxError.
			// Otherwise, the events will be unpredictable
			exp.v2 = exp.v2[1 : len(exp.v2)-1] // remove quotes
		} else if sequence == "true" || sequence == "false" {
			// Since we are overwriting the comparison with its result,
			// we need to handle such cases
			exp = expression{v1: sequence, v2: "", op: ""}
		} else {
			continue
		}

		output = append(output, exp)
	}

	return output, nil
}

// Find logical operators "and" and "or".
// Returns them in the same order.
func parseLogicalOperators(s string) []string {
	var output []string

	for _, sequence := range strings.Split(s, " ") {
		if sequence == "and" || sequence == "or" {
			output = append(output, sequence)
		}
	}

	return output
}

// Splits the string into sub-expressions and compares them to obtain the final
// comparison result.
func compareExpressionWithMarkers(s string) (bool, error) {
	comparedMarkers, err := parseMarkers(s)
	if err != nil {
		return false, err
	}

	result := comparedMarkers[0]
	for i, v := range parseLogicalOperators(s) {
		// Compare the last result with the next one to traverse the entire
		// completed expression consisting only of logical operators
		result, err = compareBool(result, comparedMarkers[i+1], v)
		if err != nil {
			return false, err
		}
	}

	return result, nil
}

// CompareMarkers parses Python markers (PEP 508) and compares them.
// Returns the comparison result or an error in case of syntax error or unknown marker.
func CompareMarkers(exp string) (bool, error) {
	for {
		startIndex, endIndex := getBracketIndexes(exp)
		if startIndex != -1 {
			// Take a slice of the expressions without brackets inside
			// the deepest pair of parentheses
			sub := exp[startIndex+1 : endIndex]
			// Compare sub-expression with markers
			result, err := compareExpressionWithMarkers(sub)
			if err != nil {
				return false, err
			}
			// Cut out the selected expression and replace it with the result
			// of the comparison
			exp = exp[:startIndex] + strconv.FormatBool(result) + exp[endIndex+1:]
		} else {
			return compareExpressionWithMarkers(exp)
		}
	}
}
