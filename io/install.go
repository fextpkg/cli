package io

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"upip/base_cfg"
	"upip/color"
	"upip/utils"
	"upip/whl"
)

type Package struct {
	Name string
	LibDir string
	Operators [][]string
}

// clean pkgName and split operators
func (pkg *Package) Prepare() {
	pkg.Name, pkg.Operators = utils.SplitOperators(pkg.Name)
}

// Returns pkgName, version, dependencies
func (pkg *Package) DefaultInstall(offset int) (string, string, []string, error) {
	pkg.Prepare()

	nameWithOffset := strings.Repeat(" ", offset * 4) + pkg.Name // used for beauty print when dependency installed
	fmt.Printf("\r%s - Scanning package..", nameWithOffset)
	version, dependencies, err := pkg.install(&Buffer{
		pkgName: pkg.Name,
		maxMessageLength: base_cfg.MAX_MESSAGE_LENGTH + len(pkg.Name),
	})

	if err != nil {
		return nameWithOffset, "", nil, err
	}

	return nameWithOffset, version, dependencies, nil
}

// Returns version, dependencies
func (pkg *Package) install(buffer interface{Write([]byte) (int, error)
	UpdateTotal(int)}) (string, []string, error) {
	doc, err := getPackageList(pkg.Name)
	if err != nil {
		return "", nil, err
	}

	pythonV := utils.ParsePythonVersion(pkg.LibDir)
	pkgVersion, link, err := selectCorrectPackageVersion(doc, pkg.Operators, pythonV)
	if err != nil {
		return "", nil, err
	}

	if CheckPackageExists(pkg.Name, pkg.LibDir, [][]string{[]string{"==", pkgVersion}}) {
		return "", nil, errors.New("Package already installed")
	}

	// uninstall all other versions of package
	UninstallPackage(pkg.LibDir, pkg.Name)

	filePath, err := downloadPackage(buffer, link, pkg.LibDir)
	if err != nil {
		return "", nil, err
	}

	err = setupPackage(filePath)
	os.RemoveAll(filePath) // remove tmp file, that was downloaded
	if err != nil {
		return "", nil, err
	}

	utils.ClearLastMessage(base_cfg.MAX_MESSAGE_LENGTH + len(pkg.Name))

	p, err := whl.LoadPackage(pkg.Name, pkg.LibDir)
	if err != nil {
		return "", nil, err
	}

	return pkgVersion, p.LoadDependencies(), nil
}

func setupPackage(pathToFile string) error {
	return whl.Extract(pathToFile)
}

// start single thread download files. Returns count downloaded packages and dependencies
func SingleThreadDownload(libDir string, packages []string, offset int) (int, int) {
	var count, dependencyCount int


	for _, name := range packages {
		pkg := Package{Name: name, LibDir: libDir}
		pkgName, pkgVersion, dependencies, err := pkg.DefaultInstall(offset)

		if err != nil {
			color.PrintflnStatusError("\r%s - Install failed", err.Error(), pkgName)
		} else {
			count++
			color.PrintflnStatusOK("\r%s (%s) - Installed", pkgName, pkgVersion)

			if len(dependencies) > 0 {
				fmt.Println(strings.Repeat(" ", offset * 4) + "-> Installing dependencies")
				c, dc := SingleThreadDownload(libDir, dependencies, offset + 1)
				dependencyCount += c + dc
			}
		}
	}

	return count, dependencyCount
}
