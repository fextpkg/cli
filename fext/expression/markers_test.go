package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ferror"
)

var (
	markerPythonVersionTrue  = "python_version >= '3." + config.GetPythonMinorVersion() + "'"
	markerPythonVersionFalse = "python_version == '3.0'"

	markerSysPlatformTrue  = "sys_platform == '" + config.SysPlatform + "'"
	markerSysPlatformFalse = "sys_platform == 'some_unknown'"

	// We can specify anything to the extra field as we don't require this marker.
	// We discover the necessary packages while parsing the metadata.
	markerExtraTrue = "extra == 'some_extra'"

	markerUnknown = "unknown == 'value'"

	markerSysPlatformUnknownOperator = "sys_platform ~= '" + config.SysPlatform + "'"

	trueMarkers = []string{
		markerPythonVersionTrue,
		markerSysPlatformTrue,
		markerExtraTrue,
	}

	falseMarkers = []string{
		markerPythonVersionFalse,
		markerSysPlatformFalse,
	}

	markersToParse = []testValue{
		{
			input:    markerSysPlatformTrue + " and true",
			expected: []bool{true, true},
		},
		{
			input:    markerSysPlatformFalse + " and false",
			expected: []bool{false, false},
		},
		{
			input:    markerExtraTrue + " and true and false and true",
			expected: []bool{true, true, false, true},
		},
	}
)

func TestCompareMarker(t *testing.T) {
	// Markers that should be true
	for _, marker := range trueMarkers {
		exp, err := parseExpressionWithOperators(marker)
		assert.Len(t, exp, 1)

		result, err := compareMarker(exp[0])
		assert.Nil(t, err)
		assert.True(t, result)
	}

	// Markers that should be false
	for _, marker := range falseMarkers {
		exp, err := parseExpressionWithOperators(marker)
		assert.Len(t, exp, 1)

		result, err := compareMarker(exp[0])
		assert.Nil(t, err)
		assert.False(t, result)
	}
}

func TestParseMarkers(t *testing.T) {
	for _, exp := range markersToParse {
		result, err := parseMarkers(exp.input)
		assert.Nil(t, err)
		assert.Equal(t, result, exp.expected)
	}
}

func TestParseMarkersInvalid(t *testing.T) {
	for _, exp := range comparisonExpressionsInvalid {
		_, err := parseMarkers(exp)
		assert.ErrorIs(t, err, ferror.SyntaxError)
	}
}

func TestUnexpectedValues(t *testing.T) {
	var unexpectedMarker *ferror.UnexpectedMarker
	var unexpectedOperator *ferror.UnexpectedOperator

	// Unexpected marker
	exp, err := parseExpressionWithOperators(markerUnknown)
	assert.Len(t, exp, 1)

	_, err = compareMarker(exp[0])
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &unexpectedMarker)

	// Unexpected operator
	exp, err = parseExpressionWithOperators(markerSysPlatformUnknownOperator)
	assert.Len(t, exp, 1)

	_, err = compareMarker(exp[0])
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &unexpectedOperator)
}

func TestMatchExtraMarker(t *testing.T) {
	match, err := MatchExtraMarker("extra == 'test'", "test")
	assert.Nil(t, err)
	assert.True(t, match)

	match, err = MatchExtraMarker("extra == 'test'", "test2")
	assert.Nil(t, err)
	assert.False(t, match)

	_, err = MatchExtraMarker("extra ==", "")
	assert.ErrorIs(t, err, ferror.SyntaxError)
}
