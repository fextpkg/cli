package cmd

import (
	"github.com/Flacy/fext/help"
	"github.com/Flacy/fext/io"
	"github.com/Flacy/fext/utils"

	"fmt"
	"io/ioutil"
	"strings"
)

func Install(libDir string, args []string) {
	options, offset := utils.ParseOptions(args)

	for _, option := range options {
		switch option {
		case "h", "help":
			help.ShowInstall()
			return
		}
	}

	packages := args[offset:]
	count, dependencyCount := io.SingleThreadDownload(libDir, packages, 0)

	fmt.Printf("\nInstalled %d packages and %d dependensies\n", count, dependencyCount)
}

func Uninstall(libDir string, args []string) {
	for _, pkgName := range args {
		io.UninstallPackage(libDir, pkgName)
	}

	fmt.Println("ok")
}

// show list of installed packages
func Freeze(path string) {
	packages, err := ioutil.ReadDir(path)

	if err != nil {
		// TODO maybe create path or notify
	}

	count := 0
	lastPkgName := ""
	// using sorted array of packages. We check each element and check if this element repeat, don't print them
	for _, pkg := range packages {
		name := pkg.Name()

		if strings.HasSuffix(name, "info") {
			// array ["name", "1.0.0", "py3..."]
			meta := strings.SplitN(name, "-", 3)
			if meta[0] != lastPkgName {
				lastPkgName = meta[0]

				fmt.Printf("%s (%s)\n", lastPkgName, utils.ClearVersion(meta[1]))
				count++
			}
		}
	}

	fmt.Printf("\nTotal: %d\n", count)
}
