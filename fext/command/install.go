package command

import (
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/io/installer"
	"github.com/fextpkg/cli/fext/ui"
)

var (
	optNoDependencies bool // without dependencies
	optSilent         bool // output only error messages
)

func getPackagesFromFiles(files []string) ([]string, error) {
	var packages []string
	for _, fileName := range files {
		p, err := io.ReadLines(fileName)
		if err != nil {
			return nil, err
		}
		packages = append(packages, p...)
	}
	return packages, nil
}

func Install(packages []string) {
	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			ui.PrintHelpInstall()
			return
		case "n", "no-dependencies":
			optNoDependencies = true
		case "s", "silent":
			optSilent = true
		case "r", "requirements":
			var err error
			packages, err = getPackagesFromFiles(packages)
			if err != nil {
				ui.PrintlnError("Failed to read file:", err.Error())
			}
		default:
			ui.PrintUnknownOption(f, ui.PrintHelpInstall)
			return
		}
	}

	i := installer.NewInstaller()
	err := i.InitializePackages(packages)
	if err != nil {
		ui.PrintlnError(err.Error())
		return
	}
	i.Install()
}
