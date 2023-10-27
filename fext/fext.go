package main

import (
	"errors"
	"os"

	"github.com/fextpkg/cli/fext/command"
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/ui"
)

type ICommand interface {
	DetectFlags() error
	Execute()
}

// getCommandInterface separates the command from arguments and returns the
// necessary command struct in accordance with the ICommand interface.
// Additionally, returns a function that outputs a docstring about specified
// command.
// Returns ferror.UnexpectedCommand if command was not found.
func getCommandInterface() (ICommand, func(), error) {
	commandName := config.Command[0]
	args := config.Command[1:]

	switch commandName {
	case "install", "i":
		return command.InitInstall(args), ui.PrintHelpInstall, nil
	case "uninstall", "u":
		return command.InitUninstall(args), ui.PrintHelpUninstall, nil
	case "freeze", "f":
		return command.InitFreeze(), ui.PrintHelpFreeze, nil
	case "debug":
		// The "debug" command doesn't accept any flags. Respectively, the "DetectFlags"
		// method will never return a "ferror.HelpFlag" error, which means that helpFunc
		// will never be used.
		return command.InitDebug(), nil, nil
	default:
		return nil, nil, ferror.UnexpectedCommand
	}
}

// Executes the selected command and its arguments, with specified options and
// flags. Prints error and terminates the process with status code 1, in case
// something went wrong.
func executeCommand() {
	cmd, helpFunc, err := getCommandInterface()
	if err != nil {
		ui.Fatal("Unable to get command:", err.Error())
	}

	if err = cmd.DetectFlags(); err != nil {
		var unknownFlag *ferror.UnknownFlag
		if errors.Is(err, ferror.HelpFlag) { // Command docstring requires
			helpFunc()
		} else if errors.As(err, &unknownFlag) { // Unknown flag
			ui.PrintUnknownOption(err.Error())
			helpFunc()
		} else { // Unexpected error
			ui.PrintlnError("Unable to detect flags:", err.Error())
		}

		os.Exit(1)
	}

	cmd.Execute()
}
func main() {
	if len(config.Command) == 0 { // No arguments were passed
		ui.PrintHelp()
	} else {
		executeCommand()
	}
}
