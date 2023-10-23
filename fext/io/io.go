package io

import (
	"bufio"
	"github.com/fextpkg/cli/fext/config"
	"os"
	"path/filepath"
	"strings"
)

// ReadLines reads lines from the file and splits them through "\n".
// Automatically trims all empty spaces
func ReadLines(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(f)
	var result []string

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line != "" {
			result = append(result, line)
		}
	}

	return result, f.Close()
}

// CreateInstallerFile creates a file named "INSTALLER" (PEP 627) in the
// specified path and writes the "fext" name there.
func CreateInstallerFile(path string) error {
	// https://peps.python.org/pep-0627/#optional-installer-file
	f, err := os.Create(filepath.Join(path, "INSTALLER"))
	if err != nil {
		return err
	}

	err = f.Chmod(config.DefaultChmod)
	if err != nil {
		return err
	}

	_, err = f.WriteString("fext\n")
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}
