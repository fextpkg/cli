package utils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

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

func parsePreVersion(s string) (int, int, error) {
	var patchValue, preValue int
	var err error
	for i, v := range s {
		if !unicode.IsDigit(v) { // find first character
			patchValue, err = strconv.Atoi(s[:i]) // cut part with characters and convert
			if err != nil {                       // unknown error
				return 0, 0, err
			}
			preValue, err = convertPreToInt(s[i:]) // convert part with characters
			if err != nil {
				return 0, 0, err
			}
			break
		}
	}
	return patchValue, preValue, nil
}

func convertPreToInt(s string) (int, error) {
	var output int
	for _, v := range s {
		output += int(v)
	}
	return output, nil
}

func compareVersion(a, b string) (int, error) {
	v1, v1pre, err := parseVersion(a)
	if err != nil {
		return 0, err
	}
	v2, v2pre, err := parseVersion(b)
	if err != nil {
		return 0, err
	}
	for i := 0; i < 2; i++ { // compare major, minor and patch version
		if v1[i] > v2[i] {
			return 1, nil
		} else if v2[i] > v1[i] {
			return -1, nil
		}
	}

	if v1pre > v2pre {
		if v2pre == 0 {
			return -1, nil
		}
		return 1, nil
	} else if v2pre > v1pre {
		if v1pre == 0 {
			return 1, nil
		}
		return -1, nil
	}

	return 0, nil
}

func CompareVersion(v1, op, v2 string) (bool, error) {
	res, err := compareVersion(v1, v2)
	if err != nil {
		return false, err
	}
	if res < 0 {
		if op == "<" || op == "<=" {
			return true, nil
		}
	} else if res > 0 {
		if op == ">" || op == ">=" {
			return true, nil
		}
	} else {
		if op == "==" || op == ">=" || op == "<=" {
			return true, nil
		}
	}
	return false, nil
}

// Split name, operator and version. Returns [][]string{operator, version} if operator was split successful
func SplitOperators(name string) (string, [][]string) {
	var operators [][]string
	// parse operators and split them. E.g. "name<=4.0.0 >=4.0.0" => [[<=, 4.0.0], [>=, 4.0.0]]
	// NOTE: separator can be anything and it also may not exists
	re, _ := regexp.Compile(`([<>!=]=?)([\d\w\.]+)`)
	v := re.FindAllStringSubmatch(strings.ReplaceAll(name, " ", ""), -1)

	for _, value := range v {
		operators = append(operators, value[1:]) // [baseValue, operator, version]
	}

	// split name
	re, _ = regexp.Compile(`[\w|\-]+`)
	name = re.FindString(name)

	return name, operators
}
