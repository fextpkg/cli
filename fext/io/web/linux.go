//go:build linux

package web

import (
	"strings"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
)

func checkCompatibility(platformTag string) (bool, error) {
	platforms := strings.Split(platformTag, ".")
	for _, platform := range platforms {
		if strings.HasPrefix(platform, "manylinux_") {
			data := strings.SplitN(platform, "_", 4)
			if !compareArch(data[3]) {
				continue
			}
			version := strings.Join(data[1:3], ".")
			res, err := expression.CompareVersion(config.GLibCVersion, ">=", version)
			if err != nil {
				return false, err
			} else if res {
				return true, nil
			}
		}
	}
	return false, nil
}

func compareArch(arch string) bool {
	return arch == config.PythonArch
}
