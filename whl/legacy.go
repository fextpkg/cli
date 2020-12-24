package whl

import (
	"strings"
)

// NOTE: This file support only READ egg format, and doesn't anything more

func loadEggDependencies(dir string) []string {
	// TODO compare sys_platform and python_version
	var dependencies []string
	rawContent, err := loadPackageContent(dir,"requires.txt")
	if err != nil {
		return dependencies
	}

	for _, v := range strings.Split(rawContent, "\n") {
		if v != "" && v[0] != 91 { // first letter != "["
			dependencies = append(dependencies, v)
		} else {
			break
		}
	}

	return dependencies
}
