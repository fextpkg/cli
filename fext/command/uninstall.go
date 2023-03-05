package command

import (
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

var (
	optCollectDep bool
)

func uninstall(pkgName string) error {
	p, err := pkg.Load(pkgName)
	if err != nil {
		return err
	}

	if optCollectDep {
		for _, dep := range p.Dependencies {
			uninstall(dep.Name)
		}
	}
	return p.Uninstall()
}

func Uninstall(packages []string) {
	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			ui.PrintHelpUninstall()
			return
		case "w", "with-dependencies":
			optCollectDep = true
		default:
			ui.PrintUnknownOption(f, ui.PrintHelpUninstall)
			return
		}
	}

	for _, pkgName := range packages {
		err := uninstall(pkgName)
		if err != nil {
			ui.PrintfError("Uninstall %s: %v\n", pkgName, err)
		}
	}
}
