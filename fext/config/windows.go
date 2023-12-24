//go:build windows

package config

import (
	"fmt"
	"os"
)

const (
	MarkerPlatform     = "win32"   // sys_platform (sys.platform)
	MarkerPlatformName = "Windows" // platform_system (platform.system())
	MarkerArch         = "AMD64"   // platform_machine (platform.machine())

	pythonExec = "python"
)

func getPythonLib() string {
	pathToAppData, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	// Python directory contains only the minor version
	return fmt.Sprintf("%s\\Python\\Python3%s\\site-packages\\", pathToAppData, GetPythonMinorVersion())
}

func getPythonVenvLib() string {
	return virtualEnvPath + "\\Lib\\site-packages\\"
}
