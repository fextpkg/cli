package command

import (
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ui"
)

func Install(packages []string) {
	//opt := io.Options{}
	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			ui.PrintHelpInstall()
			return
		case "s", "single":
			//opt.Single = true
		default:
			ui.PrintUnknownOption(f, ui.PrintHelpInstall)
			return
		}
	}
}
