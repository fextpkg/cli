package config

import (
	"log"
	"os"
	"os/exec"
)

const (
	Version      = "0.1.0"
	DefaultChmod = 0775
)

var (
	PythonVersion = getPythonVersion()
	PythonLibPath = getPythonLib() // Path to python packages directory

	virtualEnv bool

	Command []string // Command and arguments specified by user
	Flags   []string // Flags specified by user
)

func getPythonVersion() string {
	output, err := exec.Command(pythonExec, "--version").Output()
	if err != nil {
		log.Fatal("Unable to get python version. Does python exists?")
	}
	return string(output[7 : len(output)-2]) // cut off the word "Python" and the last two special characters "\r\n"
}

func checkVirtualEnv() string {
	return os.Getenv("VIRTUAL_ENV")
}

// cutQueryString used for removing all invalid dash symbols and
// to determine what given string is: part of command or flag.
func cutQueryString(string string) (string, bool) {
	for i, c := range string {
		if c != '-' {
			return string[i:], i == 0
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
	venvPath := checkVirtualEnv()
	if venvPath != "" { // venv enabled
		virtualEnv = true
	}

	if _, err := os.Stat(PythonLibPath); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(PythonLibPath, DefaultChmod); err != nil {
			log.Fatal(err)
		}
	}

	Command, Flags = parseArguments(os.Args[1:]) // First argument is a name of executable file
}
