package io

import (
	"github.com/fextpkg/cli/fext/cfg"
	"github.com/fextpkg/cli/fext/color"
	"github.com/fextpkg/cli/fext/utils"
	"github.com/fextpkg/cli/fext/whl"
	"strings"

	"errors"
	"fmt"
	"os"
)

type Options struct {
	Single bool // without dependencies
}

type Package struct {
	Name string
	Operators [][]string
	Version string
	Dependencies []string
	Size int64 // in bytes
}

// clean pkgName and split operators
func (pkg *Package) Prepare() {
	pkg.Name, pkg.Operators = utils.SplitOperators(pkg.Name)
}

// Returns pkgName
func (pkg *Package) DefaultInstall(offset int) (string, error) {
	pkg.Prepare()

	nameWithOffset := utils.GetOffsetString(offset) + pkg.Name // used for beauty print when dependency installed
	fmt.Printf("\r%s - Scanning package..", nameWithOffset)

	maxMessageLength := cfg.MAX_MESSAGE_LENGTH + len(pkg.Name)
	var err error
	err = pkg.install(&Buffer{
		pkgName:          pkg.Name,
		maxMessageLength: maxMessageLength,
	})

	utils.ClearLastMessage(maxMessageLength)
	if err != nil {
		return nameWithOffset, err
	}

	return nameWithOffset, nil
}

func (pkg *Package) install(buffer interface {
	Write([]byte) (int, error)
	UpdateTotal(int)
}) error {
	doc, err := getPackageList(pkg.Name)
	if err != nil {
		return err
	}

	var link string
	pkg.Version, link, err = selectCorrectPackageVersion(doc, pkg.Operators)
	if err != nil {
		return err
	}

	if CheckPackageExists(pkg.Name, [][]string{{"==", pkg.Version}}) {
		return errors.New("Package already installed")
	}

	// uninstall all other versions of package
	_pkg, err := whl.LoadPackage(pkg.Name)
	if err == nil {
		_pkg.Uninstall()
	}

	filePath, err := downloadPackage(buffer, link)
	if err != nil {
		return err
	}

	err = setupPackage(filePath)
	os.RemoveAll(filePath) // remove tmp file, that was downloaded
	if err != nil {
		return err
	}

	p, err := whl.LoadPackage(pkg.Name)
	if err != nil {
		return err
	}
	pkg.Size = p.GetSize()
	pkg.Dependencies, err = p.GetDependencies()
	if err != nil {
		return err
	}

	return nil
}

func setupPackage(pathToFile string) error {
	return whl.Extract(pathToFile)
}

func GetExtra(s string) (string, []string, error) {
	return getExtraNames(s)
}

// Parse extra names. Returns name, extraNames and error if syntax is invalid
func getExtraNames(s string) (string, []string, error) {
	var startQuote, endQuote int
	for i, v := range s {
		// find first quote
		if v == 91 && startQuote == 0 { // v == "["
			startQuote = i
		// find last quote
		} else if v == 93 { // v == "]"
			endQuote = i
		}
	}

	if startQuote != 0 && endQuote != 0 {
		originalName := s[:startQuote]
		s = s[startQuote+1:endQuote]
		if strings.ContainsAny(s, "[]") {
			return originalName, nil, errors.New("Syntax error")
		}

		var extraNames []string
		for _, name := range strings.Split(s, ",") {
			name = strings.ReplaceAll(name, " ", "")
			extraNames = append(extraNames, name)
		}

		return originalName, extraNames, nil
	} else if startQuote != 0 || endQuote != 0 {
		return s, nil, errors.New("Syntax error")
	}

	return s, nil, nil
}

// start single thread download files. Returns count of downloaded packages,
// dependencies and size in bytes
func SingleThreadDownload(packages []string, offset int, options *Options) (int, int, int64) {
	var count, dependencyCount int
	var size int64

	for _, name := range packages {
		name, extraNames, err := getExtraNames(name)
		if err != nil {
			color.PrintflnStatusError("%s - Unable to parse extra names", err.Error(), name)
		} else if len(extraNames) > 0 {
			extraPackages, err := GetExtraPackages(name, extraNames)
			if err != nil {
				color.PrintflnStatusError("%s - Unable to get extra names", err.Error(), name)
			} else if len(extraPackages) > 0 {
				c, dc, s := SingleThreadDownload(extraPackages, 0, options)
				count += c
				dependencyCount += dc
				size += s
			}
		} else {
			pkg := Package{Name: name}
			pkgName, err := pkg.DefaultInstall(offset)
			if err != nil {
				color.PrintflnStatusError("\r%s - Install failed", err.Error(), pkgName)
			} else {
				count++
				size += pkg.Size
				color.PrintflnStatusOK("\r%s (%s) - Installed", pkgName, pkg.Version)

				if !options.Single && len(pkg.Dependencies) > 0 {
					fmt.Println(utils.GetOffsetString(offset) + "-> Installing dependencies")
					c, dc, s := SingleThreadDownload(pkg.Dependencies, offset+1, options)
					dependencyCount += c + dc
					size += s
				}
			}
		}
	}

	return count, dependencyCount, size
}
