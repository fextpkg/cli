package whl

import (
	"github.com/Flacy/fext/fext/cfg"
	"github.com/Flacy/fext/fext/utils"

	"io/ioutil"
	"strings"
)

// Find all package directories and select wheel (dist-info) if exists, otherwise select egg-info (legacy)
func findOptimalPackageMetaDir(pkgName string) (string, error) {
	dirs := utils.GetAllPackageDirs(pkgName)

	var optimalDir string
	for _, dir := range dirs {
		optimalDir = dir

		if utils.ParseFormat(dir) == FORMAT_WHEEL {
			break
		}
	}
	return optimalDir, nil
}

func loadPackageContent(dir, fileName string) (string, error) {
	content, err := ioutil.ReadFile(cfg.PathToLib + dir + cfg.PATH_SEPARATOR + fileName)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// split description and meta, then drop part with description
func splitMeta(rawContent string) []string {
	rawContent = strings.SplitN(rawContent, "\n\n", 2)[0]
	return strings.Split(rawContent, "\n")
}

// load and parse wheel package info
func loadMeta(dir string) (*map[string]string, error) {
	pkgInfo := map[string]string{}
	rawContent, err := loadPackageContent(dir, "METADATA")
	if err != nil {
		return &pkgInfo, err
	}

	for _, v := range splitMeta(rawContent) {
		s := strings.SplitN(v, ": ", 2)
		pkgInfo[s[0]] = s[1] // key = value
	}
	delete(pkgInfo, "Description")

	return &pkgInfo, nil
}

func loadDependencies(dir string) ([]string, error) {
	var dependencies []string
	rawContent, err := loadPackageContent(dir, "METADATA")
	if err != nil {
		return nil, err
	}

	var findKey bool
	for _, v := range splitMeta(rawContent) {
		s := strings.SplitN(v, ": ", 2) // [key, value]
		if s[0] == "Requires-Dist" {
			findKey = true
			s = strings.Split(s[1], " ; ")
			if len(s) > 1 { // check if contains expression
				if ok, _ := utils.CompareExpression(s[1], dir); !ok {
					continue
				}
			}
			dependencies = append(dependencies, s[0])
		} else if findKey || s[0] == "Provides-Extra" {
			break
		}
	}

	return dependencies, nil
}

