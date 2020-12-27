package utils

import (
	"runtime"
)

// Python env markers (PEP 508)

// python_version
func GetPythonVersion(libDir string) string {
	return ParsePythonVersion(libDir)
}

// sys_platform
func GetSysPlatform() string {
	const OS = runtime.GOOS

	if OS == "windows" {
		return "win32"
	} else if OS == "linux" {
		// Please NOTE! UPIP DOESN'T support python 2, respectively value "linux2" also
		return "linux"
	} else if OS == "darwin" {
		return "darwin"
	}

	return ""
}
