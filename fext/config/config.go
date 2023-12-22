package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fextpkg/cli/fext/ui"
)

const (
	Version      = "0.4.0"
	DefaultChmod = 0755
)

var (
	virtualEnvPath = getVirtualEnvPath()

	PythonVersion string
	PythonLibPath string // Path to python packages directory

	Command []string // Command and arguments specified by user
	Flags   []string // Flags specified by user
)

func GetPythonMinorVersion() string {
	return strings.Split(PythonVersion, ".")[1]
}

func getPythonVersion() string {
	output, err := exec.Command(pythonExec, "--version").Output()
	if err != nil {
		ui.Fatal("Unable to get python version. Does python exists?")
	}

	// Cut off the word "Python" and any escaped characters
	return strings.TrimSpace(string(output[7:]))
}

func getVirtualEnvPath() string {
	return os.Getenv("VIRTUAL_ENV")
}

// cutQueryString used for removing all invalid dash symbols and
// to determine what given string is: part of command or flag.
func cutQueryString(s string) (string, bool) {
	for i, c := range s {
		if c != '-' {
			return s[i:], i == 0
		}
	}

	return "", true
}

// parseArguments is a function for parsing a user's query.
// Returns both slice with command and slice with all flags.
func parseArguments(args []string) ([]string, []string) {
	var flags, command []string
	for _, v := range args {
		cutString, isCommand := cutQueryString(v)
		if isCommand {
			command = append(command, cutString)
		} else {
			flags = append(flags, cutString)
		}
	}

	return command, flags
}

func init() {
	PythonVersion = getPythonVersion()

	// Fill in the variables based on whether the virtual environment is enabled
	if virtualEnvPath != "" {
		PythonLibPath = filepath.Clean(getPythonVenvLib())
	} else {
		PythonLibPath = filepath.Clean(getPythonLib())
	}

	// Check the presence of python library directory in the system. If not exits,
	// then try to create
	if _, err := os.Stat(PythonLibPath); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(PythonLibPath, DefaultChmod); err != nil {
			ui.Fatal(err.Error())
		}
	}

	Command, Flags = parseArguments(os.Args[1:]) // The first argument is a name of executable file
}
