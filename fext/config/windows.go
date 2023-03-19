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
	return "python"
}

func getPythonVenvExec() string {
	return virtualEnvPath + "\\Scripts\\python"
}

func getPythonLib() string {
	pathToAppData, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	// Python directory contains only minor version
	return fmt.Sprintf("%s\\Python\\Python3%s\\site-packages\\", pathToAppData, getPythonMinorVersion())
}

func getPythonVenvLib() string {
	return virtualEnvPath + "\\Lib\\site-packages\\"
}
