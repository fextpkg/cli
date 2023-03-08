//go:build linux

package config

// #include <gnu/libc-version.h>
import "C"
import (
	"fmt"
	"os"
)

const (
	pythonExec  = "python3"
	SysPlatform = "linux"
)

var GLibCVersion = C.GoString(C.gnu_get_libc_version())

func getPythonLib() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// Python directory contains only minor version
	return fmt.Sprintf("%s/.local/lib/python3.%s/site-packages/", homePath, getPythonMinorVersion())
}
