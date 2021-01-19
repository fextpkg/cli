package utils

import "github.com/fextpkg/cli/fext/cfg"

// Python env markers (PEP 508)

// python_version
func GetPythonVersion(libDir string) string {
	return ParsePythonVersion(libDir)
}

// sys_platform
func GetSysPlatform() string {
	return cfg.PYOS
}
