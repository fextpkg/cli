package cmd

import (
	"fmt"
	"github.com/Flacy/fext/help"
	"github.com/Flacy/fext/io"
	"github.com/Flacy/fext/utils"
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

	fmt.Printf("\nInstalled %d packages and %d dependencies\n", count, dependencyCount)
}

func Uninstall(libDir string, args []string) {
	options, offset := utils.ParseOptions(args)
	var collectDependency bool

	for _, option := range options {
		switch option {
		case "h", "help":
			help.ShowUninstall()
			return
		case "w", "with-dependencies":
			collectDependency = true
		}
	}

	packages := args[offset:]
	count, depCount, size := io.UninstallPackages(libDir, packages, collectDependency, false)

	fmt.Printf("\nRemoved %d packages and %d dependencies of a size %.2f MiB\n", count, depCount, float32(size / 1024) / 1024)
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
			pkgName, v, _ := utils.ParseDirectoryName(name)
			if pkgName != lastPkgName {
				lastPkgName = pkgName

				fmt.Printf("%s (%s)\n", lastPkgName, utils.ClearVersion(v))
				count++
			}
		}
	}

	fmt.Printf("\nTotal: %d\n", count)
}
