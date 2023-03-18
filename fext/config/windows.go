//go:build windows

package config

import (
	"fmt"
	"os"
)

const (
	SysPlatform = "win32" // PEP 508 marker
)

func getPythonExec() string {
	if virtualEnvPath != "" {
		return virtualEnvPath + "\\Scripts\\python"
	} else {
		return "python"
	}
}

func getPythonLib() string {
	pathToAppData, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	if virtualEnvPath != "" {
		return virtualEnvPath + "\\Lib\\site-packages\\"
	} else {
		// Python directory contains only minor version
		return fmt.Sprintf("%s\\Python\\Python3%s\\site-packages\\", pathToAppData, getPythonMinorVersion())
	}
}
