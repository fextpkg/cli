//go:build windows
// +build windows

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
	// FIXME: 3.7.0 => 7. || 3.10.0 => 10 (remove dot)
	// Trim the first 2 characters and leave 2 after them (3.10.0 => 10). Python directory contains only major version
	return fmt.Sprintf("%s\\Python\\Python3%s\\site-packages\\", pathToAppData, PythonVersion[2:][:2])
}
