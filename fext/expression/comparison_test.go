package expression

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fextpkg/cli/fext/ferror"
)

type testValue struct {
	input    string
	expected any
}

var (
	parenthesesValues = []testValue{
		{input: "()", expected: [2]int{0, 1}},
		{input: "((()))", expected: [2]int{2, 3}},
		{input: "((((", expected: [2]int{-1, -1}},
		{input: "((()", expected: [2]int{2, 3}},
	}
	comparisonExpressions = []testValue{
		{
			input: "test >= '1.0'",
			expected: []expression{
				{v1: "test", v2: "1.0", op: ">="},
			},
		},
		{
			input: "test >= '1.0' and test2 != '3' or test3 > 'some_var'",
			expected: []expression{
				{v1: "test", v2: "1.0", op: ">="},
				{v1: "test2", v2: "3", op: "!="},
				{v1: "test3", v2: "some_var", op: ">"},
			},
		},
		{
			input: "true and true",
			expected: []expression{
				{v1: "true", v2: "", op: ""},
				{v1: "true", v2: "", op: ""},
			},
		},
		{
			input: "true",
			expected: []expression{
				{v1: "true", v2: "", op: ""},
			},
		},
		{input: "", expected: []expression(nil)},
	}
	comparisonExpressionsInvalid = []string{
		"test >=",
		"<= '3.0'",
		">=",
		"test >= '3.0' and test >=",
	}

	logicalOperators = []testValue{
		{input: "true and true", expected: []string{"and"}},
		{input: "true and false or true", expected: []string{"and", "or"}},
	}

	expressionsWithMarkers = []testValue{
		{
			input: fmt.Sprintf(
				"%s and %s",
				markerPythonVersionTrue,
				markerSysPlatformTrue,
			),
			expected: true,
		},
		{
			input: fmt.Sprintf(
				"%s and %s",
				markerPythonVersionFalse,
				markerSysPlatformTrue,
			),
			expected: false,
		},
		{
			input: fmt.Sprintf(
				"%s or %s",
				markerPythonVersionFalse,
				markerSysPlatformTrue,
			),
			expected: true,
		},
		{
			input: fmt.Sprintf(
				"%s and %s and %s",
				markerPythonVersionFalse,
				markerSysPlatformTrue,
				markerSysPlatformFalse,
			),
			expected: false,
		},
		{input: "true and true and true", expected: true},
		{input: "true and false", expected: false},
		{input: markerSysPlatformFalse + " and true", expected: false},
		{input: markerSysPlatformTrue + " and true", expected: true},
		{input: markerExtraTrue + " or false", expected: true},
		{input: markerPythonVersionTrue + " and false", expected: false},
	}
	expressionsWithMarkersInvalid = []string{
		"unexpected_marker >= '1.0'",
		fmt.Sprintf("%s and unexpected_marker == 'test'", markerPythonVersionTrue),
	}

	normalExpressions = []testValue{
		{
			input: fmt.Sprintf(
				"(%s and %s) and %s",
				markerPythonVersionTrue,
				markerSysPlatformTrue,
				markerExtraTrue,
			),
			expected: true,
		},
		{
			input: fmt.Sprintf(
				"(%s and %s) or %s",
				markerPythonVersionFalse,
				markerSysPlatformTrue,
				markerExtraTrue,
			),
			expected: true,
		},
		{
			input: fmt.Sprintf(
				"(%s or %s) and %s",
				markerPythonVersionFalse,
				markerSysPlatformTrue,
				markerExtraTrue,
			),
			expected: true,
		},
		{
			input: fmt.Sprintf(
				"(%s or %s) and %s",
				markerPythonVersionFalse,
				markerSysPlatformTrue,
				markerPythonVersionFalse,
			),
			expected: false,
		},
		{
			input: fmt.Sprintf(
				"((((%s or %s)))) and (%s or %s)",
				markerPythonVersionFalse,
				markerSysPlatformTrue,
				markerPythonVersionFalse,
				markerPythonVersionTrue,
			),
			expected: true,
		},
	}
)

func TestCompareString(t *testing.T) {
	result, err := CompareString("test", "==", "test")
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = CompareString("test", "!=", "value")
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = CompareString("test", "!=", "test")
	assert.Nil(t, err)
	assert.False(t, result)

	result, err = CompareString("test", "==", "value")
	assert.Nil(t, err)
	assert.False(t, result)

	var unexpectedOperator *ferror.UnexpectedOperator
	_, err = CompareString("1", ">", "1")
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &unexpectedOperator)
}

func TestCompareBool(t *testing.T) {
	result, err := compareBool(true, true, "and")
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = compareBool(true, false, "and")
	assert.Nil(t, err)
	assert.False(t, result)

	result, err = compareBool(true, true, "or")
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = compareBool(true, false, "or")
	assert.Nil(t, err)
	assert.True(t, result)

	result, err = compareBool(false, false, "or")
	assert.Nil(t, err)
	assert.False(t, result)

	var unexpectedOperator *ferror.UnexpectedOperator
	_, err = compareBool(true, true, ">")
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &unexpectedOperator)
}

func TestGetBracketIndexes(t *testing.T) {
	for _, value := range parenthesesValues {
		startIndex, endIndex := getBracketIndexes(value.input)
		assert.EqualValues(t, [2]int{startIndex, endIndex}, value.expected)
	}
}

func TestParseExpressionWithOperators(t *testing.T) {
	for _, s := range comparisonExpressions {
		op, err := parseExpressionWithOperators(s.input)
		assert.Nil(t, err)
		assert.Equal(t, op, s.expected)
	}
}

func TestParseExpressionWithOperatorsInvalid(t *testing.T) {
	for _, s := range comparisonExpressionsInvalid {
		_, err := parseExpressionWithOperators(s)
		assert.ErrorIs(t, err, ferror.SyntaxError)
	}
}

func TestParseLogicalOperators(t *testing.T) {
	for _, value := range logicalOperators {
		result := parseLogicalOperators(value.input)
		assert.Equal(t, result, value.expected)
	}

	// Should be empty as the operators are invalid
	result := parseLogicalOperators("annd orr")
	assert.Empty(t, result)
}

func TestCompareExpressionWithMarkers(t *testing.T) {
	for _, exp := range expressionsWithMarkers {
		result, err := compareExpressionWithMarkers(exp.input)
		assert.Nil(t, err)
		assert.Equal(t, result, exp.expected)
	}
}

func TestCompareExpressionWithMarkersInvalid(t *testing.T) {
	var unexpectedMarker *ferror.UnexpectedMarker

	for _, exp := range expressionsWithMarkersInvalid {
		_, err := compareExpressionWithMarkers(exp)
		assert.ErrorAs(t, err, &unexpectedMarker)
	}
}

func TestCompareMarkers(t *testing.T) {
	for _, exp := range normalExpressions {
		result, err := CompareMarkers(exp.input)
		assert.Nil(t, err)
		assert.Equal(t, result, exp.expected)
	}
}

func TestCompareMarkersUnexpected(t *testing.T) {
	var unexpectedMarker *ferror.UnexpectedMarker

	_, err := CompareMarkers(markerUnknown)
	assert.ErrorAs(t, err, &unexpectedMarker)

	_, err = CompareMarkers(fmt.Sprintf(
		"(%s and %s) and %s",
		markerUnknown,
		markerPythonVersionTrue,
		markerSysPlatformTrue,
	),
	)
	assert.ErrorAs(t, err, &unexpectedMarker)
}
