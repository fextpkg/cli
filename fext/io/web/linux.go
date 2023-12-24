//go:build linux

package web

import (
	"strings"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
)

func checkCompatibility(platformTag string) (bool, error) {
	// There can be several platforms and they alternate through a point
	platforms := strings.Split(platformTag, ".")

	for _, platform := range platforms {
		// Select only manylinux with the glibc version,
		// since the rest are legacy
		// https://github.com/pypa/manylinux
		if strings.HasPrefix(platform, "manylinux_") {
			res, err := parseAndCompare(platform)
			if err != nil {
				return false, err
			} else if res {
				return true, nil
			}
		}
	}

	return false, nil
}

// parseAndCompare parses and compares the architecture of the platform, and also
// checks the glibc version to ensure that it is >= to the current version.
func parseAndCompare(platform string) (bool, error) {
	data := strings.SplitN(platform, "_", 4) // ["manylinux", glibc_major, glibc_minor, arch]
	if config.PythonArch != data[3] {
		return false, nil
	}

	version := strings.Join(data[1:3], ".")
	return expression.CompareVersion(config.GLibCVersion, ">=", version)
}
