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
	return "python3"
}

func getPythonVenvExec() string {
	return virtualEnvPath + "/bin/python"
}

func getPythonLib() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// Python directory contains only minor version
	return fmt.Sprintf("%s/.local/lib/python3.%s/site-packages/", homePath, getPythonMinorVersion())
}

func getPythonVenvLib() string {
	return fmt.Sprintf("%s/lib/python3.%s/site-packages/", virtualEnvPath, getPythonMinorVersion())
}