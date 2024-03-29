//go:build linux

package config

// #include <gnu/libc-version.h>
import "C"
import (
	"fmt"
	"os"
)

const (
	MarkerPlatform     = "linux" // sys_platform (sys.platform)
	MarkerPlatformName = "Linux" // platform_system (platform.system())

	pythonExec = "python3"
)

var GLibCVersion = C.GoString(C.gnu_get_libc_version())

func getPythonLib() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// Python directory contains only the minor version
	return fmt.Sprintf("%s/.local/lib/python3.%s/site-packages/", homePath, GetPythonMinorVersion())
}

func getPythonVenvLib() string {
	return fmt.Sprintf("%s/lib/python3.%s/site-packages/", virtualEnvPath, GetPythonMinorVersion())
}
