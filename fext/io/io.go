package io

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/fextpkg/cli/fext/config"
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

// ReadLinesWithComments reads file using ReadLines, but ignores commented-out
// lines starting with "#".
func ReadLinesWithComments(fileName string) ([]string, error) {
	var result []string
	rawLines, err := ReadLines(fileName)
	if err != nil {
		return nil, err
	}

	for _, line := range rawLines {
		if !strings.HasPrefix(line, "#") {
			result = append(result, line)
		}
	}

	return result, nil
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

// GetMetaDirectories goes through the directory with python modules and
// packages, selects the meta-directories and returns them.
// Returns an error if the folder could not be read.
func GetMetaDirectories() ([]string, error) {
	var directories []string

	files, err := os.ReadDir(config.PythonLibPath)
	if err != nil {
		return nil, err
	}

	// Go through the files and select meta directories
	// (wheel has the "dist-info" suffix)
	for _, f := range files {
		dirName := f.Name()
		if f.IsDir() && strings.HasSuffix(dirName, "dist-info") {
			directories = append(directories, dirName)
		}
	}

	return directories, nil
}
