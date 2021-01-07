package cmd

import (
	"fmt"
	"github.com/Flacy/fext/fext/cfg"
	"github.com/Flacy/fext/fext/color"
	"github.com/Flacy/fext/fext/help"
	"github.com/Flacy/fext/fext/io"
	"github.com/Flacy/fext/fext/utils"
	"io/ioutil"
	"os"
	"strings"
)

func Install(args []string) {
	options, offset := utils.ParseOptions(args)
	opt := io.Options{}

	for _, option := range options {
		switch option {
		case "h", "help":
			help.ShowInstall()
			return
		case "s", "single":
			opt.Single = true
		}
	}

	packages := args[offset:]
	count, dependencyCount := io.SingleThreadDownload(packages, 0, &opt)

	fmt.Printf("\nInstalled %d packages and %d dependencies\n", count, dependencyCount)
}

func Uninstall(args []string) {
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
	count, depCount, size := io.UninstallPackages(packages, collectDependency, false)

	fmt.Printf("\nRemoved %d packages and %d dependencies (%.2f MB)\n", count, depCount, float32(size / 1024) / 1024)
}

// show list of installed packages
func Freeze() {
	packages, err := ioutil.ReadDir(cfg.PathToLib)

	if err != nil {
		color.PrintfError("%s", err.Error())
		os.Exit(1)
	}

	var count int
	var size int64
	var lastPkgName string
	// using sorted array of packages. We check each element and check if this element repeat, don't print them
	for _, pkg := range packages {
		name := pkg.Name()

		if strings.HasSuffix(name, "info") {
			pkgName, v, _ := utils.ParseDirectoryName(name)
			if pkgName != lastPkgName {
				lastPkgName = pkgName

				fmt.Printf("%s (%s)\n", lastPkgName, utils.ClearVersion(v))
				count++
				size += utils.GetDirSize(pkgName)
			}
		}
	}

	fmt.Printf("\nTotal: %d (%.2f MB)\n", count, float32(size / 1024) / 1024)
}
