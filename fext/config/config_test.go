package config

import (
	"io/fs"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	commands = []string{
		"test",
		"test-command",
		"test_command",
		"+TEST command",
	}

	// {inputValue, expected}
	flags = [][2]string{
		{"-flag", "flag"},
		{"--flag", "flag"},
		{"--------flag", "flag"},
		{"-flag-test", "flag-test"},
		{"-flag---test", "flag---test"},
		{"-flag=value", "flag=value"},
	}
)

func TestPythonVersion(t *testing.T) {
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+$`, PythonVersion)
	assert.True(t, matched)
}

func TestPythonMinorVersion(t *testing.T) {
	version := GetPythonMinorVersion()
	matched, _ := regexp.MatchString(`^\d+$`, version)

	assert.True(t, matched)
}

func TestPythonPath(t *testing.T) {
	assert.True(t, fs.ValidPath(PythonLibPath))
}

func TestVirtualEnvPath(t *testing.T) {
	// Tests are executed without using a virtual env
	assert.Empty(t, virtualEnvPath)

	// Manually set up the virtual env
	newPath := "/usr/lib"
	err := os.Setenv("VIRTUAL_ENV", newPath)
	if err != nil {
		panic(err)
	}
	// And verify that the changes have been applied
	assert.Equal(t, getVirtualEnvPath(), newPath)
}

func TestCutQueryString(t *testing.T) {
	for _, cmd := range commands {
		s, isCommand := cutQueryString(cmd)

		assert.True(t, isCommand)
		assert.Equal(t, s, cmd)
	}

	for _, flag := range flags {
		inputValue, expected := flag[0], flag[1]
		s, isCommand := cutQueryString(inputValue)

		assert.False(t, isCommand)
		assert.Equal(t, s, expected)
	}
}

func TestParseArguments(t *testing.T) {
	var inputFlags []string
	for _, flag := range flags {
		inputFlags = append(inputFlags, flag[0])
	}

	parsedCommands, parsedFlags := parseArguments(commands)
	assert.Len(t, parsedCommands, len(commands))
	assert.Len(t, parsedFlags, 0)

	parsedCommands, parsedFlags = parseArguments(inputFlags)
	assert.Len(t, parsedCommands, 0)
	assert.Len(t, parsedFlags, len(inputFlags))

	parsedCommands, parsedFlags = parseArguments(append(commands, inputFlags...))
	assert.Len(t, parsedCommands, len(commands))
	assert.Len(t, parsedFlags, len(inputFlags))
}
