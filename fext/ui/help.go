package ui

import (
	"fmt"
)

// PrintHelp prints main help info
func PrintHelp() {
	fmt.Println("Usage:\n\tfext [options] <command> [args]",
		"\n\nAvailable commands:\n",
		"\t(i)nstall <package(s)>   - install a package(s)\n",
		"\t(u)ninstall <package(s)> - uninstall a package(s)\n",
		"\t(f)reeze                 - show list of installed packages\n",
		"\tshow <package>           - show general info about package\n",
		"\tcheck                    - verify correct installation of packages in the system\n",
		"\tdebug                    - show debug info",
		"\n\nFor additional help you can write:\n\tfext <command> -h")
}

// PrintHelpInstall prints install help info
func PrintHelpInstall() {
	fmt.Println("Available options:\n",
		"\t-n, --no-dependencies      - Install single package, without dependencies\n",
		"\t-q, --quiet (-s, --silent) - Print only error messages\n",
		"\t-r, --requirements         - Install from files",
	)
}

func PrintHelpUninstall() {
	fmt.Println("Available options:\n",
		"\t-d, --dependencies - Remove dependencies of package also")
}

func PrintHelpFreeze() {
	fmt.Println("Available options:\n",
		"\t-m, --mode=<str> - set the print mode: human (default), pip")
}

func PrintUnknownOption(opt string) {
	PrintlnError("Unknown option:", opt)
}
