package io

import (
	"fmt"
	"github.com/Flacy/fext/fext/color"
	"github.com/Flacy/fext/fext/utils"
	"github.com/Flacy/fext/fext/whl"
)

type Buffer struct {
	pkgName string
	maxMessageLength int // need for beauty clear progress bar
	DownloadedBytes int
	Total int
}

func (b *Buffer) Write(data []byte) (int, error) {
	count := len(data)
	b.DownloadedBytes += count / 1024 // convert to KiB
	b.updateProgressBar()

	return count, nil
}

func (b *Buffer) updateProgressBar() {
	utils.ClearLastMessage(b.maxMessageLength)

	fmt.Printf("\r%s - Downloading.. (%d/%d KiB)",
				b.pkgName, b.DownloadedBytes, b.Total)
}

func (b *Buffer) UpdateTotal(value int) {
	b.Total = value / 1024
}

func CheckPackageExists(name, libDir string, operators [][]string) bool {
	dirName := utils.GetFirstPackageMetaDir(libDir, name)

	if dirName != "" {
		_, version, _ := utils.ParseDirectoryName(dirName)
		for _, op := range operators {
			if ok, err := utils.CompareVersion(version, op[0], op[1]); err != nil || !ok {
				return false
			}
		}
		return true
	}

	return false
}

// Returns count of removed packages and total removed size in MB
func UninstallPackages(libDir string, packages []string, collectDependencies, inRecurse bool, ) (int, int, int64) {
	var count, depCount int
	var size int64
	var dependencies []string
	for _, pkgName := range packages {
		pkg, err := whl.LoadPackage(pkgName, libDir)
		if inRecurse {
			// add offset for beauty print
			pkgName = utils.GetOffsetString(1) + pkgName
		}
		if err != nil {
			color.PrintflnStatusError("%s - Uninstall failed", err.Error(), pkgName)
			continue
		} else if collectDependencies {
			dependencies = pkg.LoadDependencies()
		}

		curSize := pkg.GetSize()
		err = pkg.Uninstall(libDir)
		if err != nil {
			color.PrintflnStatusError("%s - Uninstall failed", err.Error(), pkgName)
		} else {
			color.PrintflnStatusOK("%s - Uninstalled", pkgName)
			count++
			size += curSize
		}

		// we don't run the recursion, via collectDependencies arg,
		// as this could lead to the removal most part of the packets
		if len(dependencies) > 0 {
			for i, dep := range dependencies {
				// clean name
				dependencies[i], _ = utils.SplitOperators(dep)
			}
			fmt.Println("-> Uninstalling dependencies")
			c, _, s := UninstallPackages(libDir, dependencies, false, true)
			depCount += c
			size += s
		}

	}

	return count, depCount, size
}