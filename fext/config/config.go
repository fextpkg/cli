package config

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	Version      = "0.1.1"
	DefaultChmod = 0755
)

var (
	virtualEnvPath = getVirtualEnvPath()

	pythonExec    string // Path to python executor
	PythonVersion string
	PythonLibPath string // Path to python packages directory

	Command []string // Command and arguments specified by user
	Flags   []string // Flags specified by user
)

func getPythonVersion() string {
	output, err := exec.Command(pythonExec, "--version").Output()
	if err != nil {
		log.Fatal("Unable to get python version. Does python exists?")
	}
	// Cut off the word "Python". We do not clear the last characters of \r\n,
	// because during version comparing, the strconv function is used, which clears
	// them itself
	return string(output[7:])
}

func getPythonMinorVersion() string {
	return strings.Split(PythonVersion, ".")[1]
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

// parseArguments is a function for parsing user's query.
// Returns both slice with all flags and slice with command.
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
	// Fill in the variables based on whether the virtual environment is enabled
	if virtualEnvPath != "" {
		pythonExec = getPythonVenvExec()
		PythonVersion = getPythonVersion()
		PythonLibPath = getPythonVenvLib()
	} else {
		pythonExec = getPythonExec()
		PythonVersion = getPythonVersion()
		PythonLibPath = getPythonLib()
	}

	// Check the presence of python library directory in the system. if not exits,
	// then try to create
	if _, err := os.Stat(PythonLibPath); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(PythonLibPath, DefaultChmod); err != nil {
			log.Fatal(err)
		}
	}

	Command, Flags = parseArguments(os.Args[1:]) // First argument is a name of executable file
}
