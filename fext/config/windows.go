//go:build windows

package config

import (
	"fmt"
	"os"
)

const (
	pythonExec  = "python"
	SysPlatform = "win32" // PEP 508 marker
)

func getPythonLib() string {
	pathToAppData, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	// TODO venv support
	// Python directory contains only minor version
	return fmt.Sprintf("%s\\Python\\Python3%s\\site-packages\\", pathToAppData, getPythonMinorVersion())
}
