package ui

import (
	"github.com/fextpkg/cli/fext/color"
	"github.com/fextpkg/cli/fext/config"

	"fmt"
)

// PrintHelp prints main help info
func PrintHelp() {
	fmt.Println("Usage:\n\tfext <command> [args]",
		"\n\nAvailable commands:\n",
		"\t(i)nstall [options] <package(s)> - install a package(s)\n",
		"\t(u)ninstall [options] <package(s)> - uninstall a package(s)\n",
		"\tfreeze - show list of installed packages\n",
		"\tdebug - show debug info",
		"\n\nFor additional help you can write:\n\tfext <command> -h")
}

// PrintHelpInstall prints install help info
func PrintHelpInstall() {
	fmt.Println("Available options:\n",
		"\t-s, --single - Install single package, without dependencies")
}

func PrintHelpUninstall() {
	fmt.Println("Available options:\n",
		"\t-w, --with-dependencies - Remove dependencies of package also")
}

// PrintDebug prints debug info
func PrintDebug() {
	fmt.Printf("Fext (%s)\n\nLinked to: %s\n", config.Version, config.PythonLibPath)
}

func PrintUnknownOption(opt string, call func()) {
	color.PrintfError("Unknown option: %s\n", opt)
	call()
}
