package main

import (
	"fmt"

	"github.com/fextpkg/cli/fext/command"
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ui"
)

func main() {
	if len(config.Command) == 0 {
		ui.PrintHelp()
	} else {
		cmd := config.Command[0]
		args := config.Command[1:]

		switch cmd {
		case "install", "i":
			command.Install(args)
		case "uninstall", "u":
			command.Uninstall(args)
		case "freeze":
			command.Freeze()
		case "debug":
			ui.PrintDebug()
		default:
			fmt.Println("Unexpected command")
		}
	}
}
