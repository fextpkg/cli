package io

import (
	"github.com/Flacy/fext/fext/base_cfg"
	"github.com/Flacy/fext/fext/color"
	"github.com/Flacy/fext/fext/utils"
	"github.com/Flacy/fext/fext/whl"

	"errors"
	"fmt"
	"os"
)

type Package struct {
	Name      string
	LibDir    string
	Operators [][]string
}

// clean pkgName and split operators
func (pkg *Package) Prepare() {
	pkg.Name, pkg.Operators = utils.SplitOperators(pkg.Name)
}

// Returns pkgName, version, dependencies
func (pkg *Package) DefaultInstall(offset int) (string, string, []string, error) {
	pkg.Prepare()

	nameWithOffset := utils.GetOffsetString(offset) + pkg.Name // used for beauty print when dependency installed
	fmt.Printf("\r%s - Scanning package..", nameWithOffset)

	maxMessageLength := base_cfg.MAX_MESSAGE_LENGTH + len(pkg.Name)
	version, dependencies, err := pkg.install(&Buffer{
		pkgName:          pkg.Name,
		maxMessageLength: maxMessageLength,
	})

	utils.ClearLastMessage(maxMessageLength)
	if err != nil {
		return nameWithOffset, "", nil, err
	}

	return nameWithOffset, version, dependencies, nil
}

// Returns version, dependencies
func (pkg *Package) install(buffer interface {
	Write([]byte) (int, error)
	UpdateTotal(int)
}) (string, []string, error) {
	doc, err := getPackageList(pkg.Name)
	if err != nil {
		return "", nil, err
	}

	pythonV := utils.ParsePythonVersion(pkg.LibDir)
	pkgVersion, link, err := selectCorrectPackageVersion(doc, pkg.Operators, pythonV)
	if err != nil {
		return "", nil, err
	}

	if CheckPackageExists(pkg.Name, pkg.LibDir, [][]string{{"==", pkgVersion}}) {
		return "", nil, errors.New("Package already installed")
	}

	// uninstall all other versions of package
	_pkg, err := whl.LoadPackage(pkg.Name, pkg.LibDir)
	if err == nil {
		_pkg.Uninstall(pkg.LibDir)
	}

	filePath, err := downloadPackage(buffer, link, pkg.LibDir)
	if err != nil {
		return "", nil, err
	}

	err = setupPackage(filePath)
	os.RemoveAll(filePath) // remove tmp file, that was downloaded
	if err != nil {
		return "", nil, err
	}

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
				fmt.Println(utils.GetOffsetString(offset) + "-> Installing dependencies")
				c, dc := SingleThreadDownload(libDir, dependencies, offset+1)
				dependencyCount += c + dc
			}
		}
	}

	return count, dependencyCount
}
