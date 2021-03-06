package cmd

import (
	"github.com/fextpkg/cli/fext/cfg"
	"github.com/fextpkg/cli/fext/color"
	"github.com/fextpkg/cli/fext/help"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/utils"

	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func printUnknownOption(opt string, call func()) {
	color.PrintfError("Unknown option: %s\n", opt)
	call()
}

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
		default:
			printUnknownOption(option, help.ShowInstall)
			return
		}
	}

	packages := args[offset:]
	count, dependencyCount, size := io.SingleThreadDownload(packages, 0, &opt)

	fmt.Printf(
		"\nInstalled %d packages and %d dependencies (%.2f MB)\n",
		count,
		dependencyCount,
		float32(size / 1024) / 1024,
	)
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
		default:
			printUnknownOption(option, help.ShowUninstall)
			return
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
