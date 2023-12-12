package expression

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Condition struct {
	Value    string
	Operator string
}

// parseVersion splits the version into semantic parts (major, minor, patch, pre)
// and returns them in the form [3]int{major, minor, patch}, pre. Returns an error
// if version could not be converted to int
func parseVersion(s string) ([3]int, int, error) {
	var output [3]int
	var pre int
	parts := strings.Split(s, ".")
	length := len(parts)

	for i := 0; i < 3; i++ {
		if i < length {
			value := parts[i]
			intValue, err := strconv.Atoi(value)
			if err != nil { // string contains characters
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

// parsePreVersion separates patch and pre version. Example: 1a2 => 1, 147
// (a is 97, 2 is 50 => 97 + 50 = 147). Returns an error if version could not be
// converted to int
func parsePreVersion(s string) (int, int, error) {
	var patchValue, preValue int
	var err error
	for i, v := range s {
		if !unicode.IsDigit(v) { // find first character
			patchValue, err = strconv.Atoi(s[:i]) // cut part with characters and convert
			if err != nil {                       // unknown error
				return 0, 0, err
			}
			preValue, err = convertStrToInt(s[i:]) // convert part with characters
			if err != nil {
				return 0, 0, err
			}
			break
		}
	}
	return patchValue, preValue, nil
}

func convertStrToInt(s string) (int, error) {
	var output int
	for _, v := range s {
		output += int(v)
	}
	return output, nil
}

// compareVersion returns the result of a comparison between versions. If a > b,
// 1 will be returned, if a < b -1 will be returned, if a == b 0 will be returned
func compareVersion(a, b string) (int, error) {
	v1, v1pre, err := parseVersion(a)
	if err != nil {
		return 0, err
	}
	v2, v2pre, err := parseVersion(b)
	if err != nil {
		return 0, err
	}

	for i := 0; i < 3; i++ { // compare major, minor and patch version
		if v1[i] > v2[i] {
			return 1, nil
		} else if v2[i] > v1[i] {
			return -1, nil
		}
	}

	if v1pre&v2pre != 0 {
		if v1pre > v2pre {
			return 1, nil
		} else if v2pre > v1pre {
			return -1, nil
		}
	}

	return 0, nil
}

// CompareVersion works by means of comparing each part of the version. The
// version is divided into semantic parts (major, minor, patch, pre) and
// converted into numbers that are compared one after another. For example:
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
	}
	return false, nil
}

// ParseConditions split package name and conditions. Separator can be anything,
// and it also may not exist. Returns package name, conditions. Example:
// "name<=4.0.0 >=4.0.0" => name, [(4.0.0, <=), (4.0.0, >=)]
func ParseConditions(exp string) (string, []Condition) {
	var cond []Condition
	re := regexp.MustCompile(`([<>!=]=?)([\w\.]+)`)
	v := re.FindAllStringSubmatch(strings.ReplaceAll(exp, " ", ""), -1)

	for _, value := range v {
		cond = append(cond, Condition{
			Value:    value[2],
			Operator: value[1],
		}) // value[baseValue, operator, value]
	}

	// split name
	re, _ = regexp.Compile(`[\w|\-.\[\]]+`)
	return re.FindString(exp), cond
}
