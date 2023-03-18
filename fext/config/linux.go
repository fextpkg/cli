//go:build linux

package config

// #include <gnu/libc-version.h>
import "C"
import (
	"fmt"
	"os"
)

const (
	SysPlatform = "linux"
)

var GLibCVersion = C.GoString(C.gnu_get_libc_version())

func getPythonExec() string {
	if virtualEnvPath != "" {
		return virtualEnvPath + "/bin/python"
	} else {
		return "python3"
	}
}

func getPythonLib() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	if virtualEnvPath != "" {
		return fmt.Sprintf("%s/lib/python3.%s/site-packages/", virtualEnvPath, getPythonMinorVersion())
	} else {
		// Python directory contains only minor version
		return fmt.Sprintf("%s/.local/lib/python3.%s/site-packages/", homePath, getPythonMinorVersion())
	}
}
