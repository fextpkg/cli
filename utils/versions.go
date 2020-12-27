package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// Need for correct check, cause version can have two or more digits in section
func equateVersions(v1, v2 string) (a, b string) {
	v1l := len(v1)
	v2l := len(v2)

	// make strings of same length
	if v1l > v2l {
		v2 += strings.Repeat("0", v1l - v2l)
	} else if v2l > v1l {
		v1 += strings.Repeat("0", v2l - v1l)
	}

	return v1, v2
}

// Returns versions (E.g. 4.0.0a0 => 40010)
func parseAndEquateVersions(rawv1, rawv2 string) (int, int) {
	a := strings.Split(rawv1, ".")
	b := strings.Split(rawv2, ".")
	var v1, v2 string
	var s1, s2 string // tmp values
	minLength := FindMinValue([]int{len(a), len(b)}) // drop part with more length

	for i := 0; i < minLength; i++ {
		// try to convert to int, and if it success we skip part with letters
		s1 = a[i]
		s2 = b[i]
		if _, err := strconv.Atoi(s1 + s2); err != nil {
			// convert part with letters to int
			s1, s2 = convertLettersToIntString(s1), convertLettersToIntString(s2)
		}

		s1, s2 = equateVersions(s1, s2)
		v1 += s1
		v2 += s2
	}

	// skip errors cause above we check it
	v1c, _ := strconv.Atoi(v1)
	v2c, _ := strconv.Atoi(v2)

	return v1c, v2c
}

/* compare <a> <operator> <b> and return bool result
For example: ("4.0.0a", "<=", "4.0.0") = (400, "<=", 400) => true
WARNING: Letters are also integers, starts with 1 */
func CompareVersion(a, operator, b string) (bool, error) {
	v1, v2 := parseAndEquateVersions(a, b)
	return compare(v1, v2, operator)
}

// Split name, operator and version. Returns [][]string{operator, version} if operator was split successful
func SplitOperators(name string) (string, [][]string) {
	var operators [][]string
	// parse operators and split them. E.g. "name<=4.0.0 >=4.0.0" => [[<=, 4.0.0], [>=, 4.0.0]]
	// NOTE: separator can be anything and it also may not exists
	re, _ := regexp.Compile(`([<>!=]=?)([\d\w\.]+|\,+)`)
	v := re.FindAllStringSubmatch(name, -1)

	for _, value := range v {
		operators = append(operators, value[1:]) // [baseValue, operator, version]
	}

	// split name
	re, _ = regexp.Compile(`[\w|\-]+`)
	name = re.FindString(name)

	return name, operators
}

func ClearVersion(version string) string {
	return strings.Replace(version, ".dist", "", 1)
}

// Returns name, version, format
func ParseDirectoryName(dirName string) (string, string, string) {
	// array [name, version, ...]
	meta := strings.SplitN(dirName, "-", 3)
	format := ParseFormat(dirName)

	if len(meta) <= 3 {
		meta[1] = ClearVersion(meta[1])
	}
	return meta[0], meta[1], format
}

