package expression

import (
	"strconv"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ferror"
)

func compareMarkerPythonVersion(exp expression) (bool, error) {
	return CompareVersion(config.PythonVersion, exp.op, exp.v2)
}

func compareMarkerSysPlatform(exp expression) (bool, error) {
	return CompareString(config.SysPlatform, exp.op, exp.v2)
}

func compareMarkerExtra(_ expression) (bool, error) {
	return true, nil
}

// compareMarker process and compares the given marker following the PEP 508 standard.
// If there is a syntax violation, an error is returned.
func compareMarker(exp expression) (bool, error) {
	var compareFunc func(exp expression) (bool, error)

	switch exp.v1 {
	case "python_version", "python_full_version":
		compareFunc = compareMarkerPythonVersion
	case "sys_platform":
		compareFunc = compareMarkerSysPlatform
	case "extra":
		// As we handle the discovery and validation of extra in the "MatchExtraMarker"
		// function, here we can utilize a placeholder that consistently returns true
		compareFunc = compareMarkerExtra
	default:
		return false, &ferror.UnexpectedMarker{Marker: exp.v1}
	}

	return compareFunc(exp)
}

// Parses and compares python markers (PEP 508).
// Returns the comparison result of each marker as a list.
// Returns an error in case of invalid syntax or if an unknown marker is passed.
func parseMarkers(s string) ([]bool, error) {
	var compareResults []bool
	var result bool

	expressions, err := parseExpressionWithOperators(s)
	if err != nil {
		return nil, err
	}

	for _, exp := range expressions {
		if exp.v2 != "" { // Normal expression
			result, err = compareMarker(exp)
			if err != nil {
				return nil, err
			}
		} else { // Bool result from past operations
			result, err = strconv.ParseBool(exp.v1)
			if err != nil {
				return nil, err
			}
		}

		compareResults = append(compareResults, result)
	}

	return compareResults, nil
}

// MatchExtraMarker checks for the existence of the "extra" marker and matches it
// with the given "extraName" parameter.
// Returns an error in case of syntax violation.
func MatchExtraMarker(s, extraName string) (bool, error) {
	expressions, err := parseExpressionWithOperators(s)
	if err != nil {
		return false, err
	}

	for _, exp := range expressions {
		if exp.v1 == "extra" && exp.v2 == extraName {
			return true, nil
		}
	}

	return false, nil
}
