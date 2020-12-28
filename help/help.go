package help

import (
	"github.com/Flacy/upip/base_cfg"

	"fmt"
)

// print main help info
func Show() {
	fmt.Println("Usage:\n\tupip <command> [args]",
				"\n\nAvailable commands:\n",
				"\t(i)nstall [options] <package(s)> - install a package(s)\n",
				"\t(u)ninstall <package(s)> - uninstall a package(s)\n",
				"\tfreeze - show list of installed packages\n",
				"\tdebug - show debug info")
}

// print install help info
func ShowInstall() {
	fmt.Println("Available options:")
				//"\t-t, --thread - Multithread installation\n",
				//"\t-g, --global - Install packages global (avoid venv)" // TODO
}

// print debug info
func ShowDebug(pathToConfig string, libDir string) {
	fmt.Printf("Upgraded Python Indexing Package (%s)\n\nUsing config: %s\nLinked to: %s\n",
				base_cfg.VERSION, pathToConfig, libDir)
}

