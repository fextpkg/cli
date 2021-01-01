package help

import (
	"github.com/Flacy/fext/src/base_cfg"

	"fmt"
)

// print main help info
func Show() {
	fmt.Println("Usage:\n\tfext <command> [args]",
				"\n\nAvailable commands:\n",
				"\t(i)nstall [options] <package(s)> - install a package(s)\n",
				"\t(u)ninstall [options] <package(s)> - uninstall a package(s)\n",
				"\tfreeze - show list of installed packages\n",
				"\tdebug - show debug info",
				"\n\nFor additional help you can write: \n\tfext <command> -h")
}

// print install help info
func ShowInstall() {
	fmt.Println("Available options:")
				//"\t-t, --thread - Multithread installation\n",
				//"\t-g, --global - Install packages global (avoid venv)" // TODO
}

func ShowUninstall() {
	fmt.Println("Available options:\n",
				"\t-w, --with-dependencies - Remove dependencies of package also")
}

// print debug info
func ShowDebug(pathToConfig string, libDir string) {
	fmt.Printf("FEXT (%s)\n\nUsing config: %s\nLinked to: %s\n",
				base_cfg.VERSION, pathToConfig, libDir)
}

