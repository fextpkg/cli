package io

import (
	"github.com/Flacy/fext/utils"

	"fmt"
	"os"
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

func UninstallPackage(libDir, name string) {
	dirs := utils.GetAllPackageDirs(name, libDir)

	for _, dir := range dirs {
		os.RemoveAll(libDir + dir)
	}
}
