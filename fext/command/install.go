package command

import (
	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/io/installer"
	"github.com/fextpkg/cli/fext/ui"
)

type Install struct {
	// Load packages names from the list of files passed in arguments instead
	// of packages names
	fileMode bool

	// List of packages to be installed
	packages []string

	// Installation options. Are filled in based on the passed flags
	options *installer.Options
}

// getPackagesFromFiles retrieves the list of packages names from the files
// passed instead of the package list. This method is used when fileMode is
// enabled.
func (cmd *Install) getPackagesFromFiles() ([]string, error) {
	var packages []string
	for _, fileName := range cmd.packages {
		p, err := io.ReadLinesWithComments(fileName)
		if err != nil {
			return nil, err
		}
		packages = append(packages, p...)
	}

	return packages, nil
}

// Installs the list of passed packages. Returns an error if extra packages could
// not be retrievers (they do not exist or load error).
func (cmd *Install) install() error {
	i := installer.NewInstaller(cmd.options)
	err := i.InitializePackages(cmd.packages)
	if err != nil {
		return err
	}

	i.Install()
	return nil
}

// DetectFlags analyzes the passed flags and fills in the variables associated
// with them.
//
// Returns ferror.HelpFlag if you need to print the docstring about this command.
// Returns ferror.UnknownFlag if passed the unknown flag.
func (cmd *Install) DetectFlags() error {
	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			return ferror.HelpFlag
		case "n", "no-deps", "no-dependencies":
			cmd.options.NoDependencies = true
		case "s", "silent", "q", "quiet":
			cmd.options.QuietMode = true
		case "r", "requirements":
			cmd.fileMode = true
		default:
			return &ferror.UnknownFlag{Flag: f}
		}
	}

	return nil
}

// Execute downloads and installs the passed packages using the flags set.
// Additionally, scans files if fileMode is enabled.
func (cmd *Install) Execute() {
	var err error
	if cmd.fileMode {
		cmd.packages, err = cmd.getPackagesFromFiles()
		if err != nil {
			ui.Fatal("Failed to read file:", err.Error())
		}
	}

	if err = cmd.install(); err != nil {
		ui.Fatal("Unable to install:", err.Error())
	}
}

// InitInstall initializes "install" command structure with the default
// parameters. Takes as an argument a list of packages names, or filenames that
// includes package names, to be installed. Each package name can contain either
// operators or extra names.
func InitInstall(packages []string) *Install {
	return &Install{
		fileMode: false,
		packages: packages,
		options:  installer.DefaultOptions(),
	}
}
