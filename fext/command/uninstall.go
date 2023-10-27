package command

import (
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

type Uninstall struct {
	// Flags
	collectDependencies bool // Uninstall dependencies of packages

	// Internal variables
	packages []string // List of packages to be uninstalled
}

// Removes the specified package and all associated files. Also removes its
// dependencies if collectDependencies is true.
func (cmd *Uninstall) uninstall(pkgName string) error {
	p, err := pkg.Load(pkgName)
	if err != nil {
		return err
	}

	if cmd.collectDependencies {
		for _, dep := range p.Dependencies {
			// Recursion is used here because the uninstallation command is not in use
			// priority. In the future it will be redone by safe uninstalling (issue #1)
			cmd.uninstall(dep.Name)
		}
	}

	return p.Uninstall()
}

// DetectFlags analyzes the passed flags and fills in the variables associated
// with them.
//
// Returns ferror.HelpFlag if you need to print the docstring about this command.
// Returns ferror.UnknownFlag if passed the unknown flag.
func (cmd *Uninstall) DetectFlags() error {
	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			return ferror.HelpFlag
		case "d", "dependencies":
			cmd.collectDependencies = true
		default:
			return &ferror.UnknownFlag{Flag: f}
		}
	}

	return nil
}

// Execute removes the passed packages and all associated files.
func (cmd *Uninstall) Execute() {
	for _, pkgName := range cmd.packages {
		err := cmd.uninstall(pkgName)
		if err != nil {
			ui.PrintfError("Uninstall %s: %v\n", pkgName, err)
		}
	}
}

// InitUninstall initializes "uninstall" command structure with the default
// parameters. Takes as an argument a list of packages names to be deleted.
func InitUninstall(packages []string) *Uninstall {
	return &Uninstall{
		collectDependencies: false,
		packages:            packages,
	}
}
