package command

import (
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/io/installer"
	"github.com/fextpkg/cli/fext/ui"
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
	opt := installer.DefaultOptions()

	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			ui.PrintHelpInstall()
			return
		case "n", "no-dependencies":
			opt.NoDependencies = true
		case "s", "silent", "q", "quiet":
			opt.Silent = true
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

	i := installer.NewInstaller(opt)
	err := i.InitializePackages(packages)
	if err != nil {
		ui.PrintlnError("Failed to initialize packages:", err.Error())
		return
	}
	i.Install()
}
