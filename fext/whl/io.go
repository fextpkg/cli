package whl

import (
	"github.com/fextpkg/cli/fext/cfg"
	"github.com/fextpkg/cli/fext/utils"

	"fmt"
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

// Load and parse meta data. Returns [[key, value]]
func loadPackageMetaData(metaDir string) ([][]string, error) {
	content, err := ioutil.ReadFile(cfg.PathToLib + metaDir + cfg.PATH_SEPARATOR + "METADATA")
	if err != nil {
		return nil, err
	}

	var metaData [][]string
	for _, v := range strings.Split(string(content), "\n") {
		if !strings.HasSuffix(v, " ") { // skip description and other strings
			metaData = append(metaData, strings.SplitN(v, ": ", 2))
		}
	}

	return metaData, nil
}

// split description and meta, then drop part with description
func splitMeta(rawContent string) []string {
	rawContent = strings.SplitN(rawContent, "\n\n", 2)[0]
	return strings.Split(rawContent, "\n")
}

// load and parse wheel package info
func loadMeta(metaDir string) (*map[string]string, error) {
	// TODO : refactor functions for correct work in all versions
	pkgInfo := map[string]string{}

	return &pkgInfo, nil
}

// Returns dependencies, extra without the result of comparisons
func loadRawDependenciesAndExtra(metaDir string) ([]string, []string, error) {
	var dependencies, extra []string
	metaData, err := loadPackageMetaData(metaDir)
	if err != nil {
		return nil, nil, err
	}

	for _, v := range metaData {
		if v[0] == "Requires-Dist" {
			value := v[1]
			if strings.Contains(value, "extra") {
				extra = append(extra, value)
			} else {
				dependencies = append(dependencies, value)
			}
		}
	}

	return dependencies, extra, nil
}

// Parse and compare dependencies
func parseDependencies(metaDir string, rawDependencies []string) ([]string, error) {
	var dependencies []string
	for _, dep := range rawDependencies {
		v := strings.SplitN(dep, " ; ", 2)
		if len(v) > 1 { // check if value have expression
			ok, err := utils.CompareExpression(dep, metaDir)
			if err != nil {
				return nil, err
			} else if !ok {
				continue
			}
		}

		dependencies = append(dependencies, v[0])
	}

	return dependencies, nil
}

// Parse and compare raw extra and returns packages provides by extra
func parseExtra(name, metaDir string, rawExtra []string) ([]string, error) {
	var extraPackages []string
	for _, v := range rawExtra {
		data := strings.SplitN(v, " ; ", 2) // [name, expression]
		exp := fmt.Sprintf("extra == '%s'", name)
		if strings.Contains(data[1], exp) {
			data[1] = strings.ReplaceAll(data[1], exp, "true")
			ok, err := utils.CompareExpression(data[1], metaDir)
			if err != nil {
				return nil, err
			} else if ok {
				extraPackages = append(extraPackages, data[0])
			}
		}
	}

	return extraPackages, nil
}
