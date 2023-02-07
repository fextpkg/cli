package io

import (
	"github.com/fextpkg/cli/fext/color"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/utils"
	"strings"

	"errors"
	"fmt"
	"os"
)

type Options struct {
	Single bool // without dependencies
}

type Package struct {
	Name         string
	Operators    [][]string
	Version      string
	Dependencies []string
	Size         int64 // in bytes
}

// clean pkgName and split operators
func (p *Package) Prepare() {
	p.Name, p.Operators = utils.SplitOperators(p.Name)
}

// Returns pkgName
func (p *Package) DefaultInstall(offset int) (string, error) {
	p.Prepare()

	nameWithOffset := utils.GetOffsetString(offset) + p.Name // used for beauty print when dependency installed
	fmt.Printf("\r%s - Scanning pkg..", nameWithOffset)

	maxMessageLength := 20 + len(p.Name)
	var err error
	err = p.install(&Buffer{
		pkgName:          p.Name,
		maxMessageLength: maxMessageLength,
	})

	utils.ClearLastMessage(maxMessageLength)
	if err != nil {
		return nameWithOffset, err
	}

	return nameWithOffset, nil
}

func (p *Package) install(buffer interface {
	Write([]byte) (int, error)
	UpdateTotal(int)
}) error {
	doc, err := getPackageList(p.Name)
	if err != nil {
		return err
	}

	var link string
	p.Version, link, err = selectCorrectPackageVersion(doc, p.Operators)
	if err != nil {
		return err
	}

	//if CheckPackageExists(p.Name, [][]string{{"==", p.Version}}) {
	//	return errors.New("Package already installed")
	//}
	// TODO move check above to the bottom lines of code

	// uninstall all other versions of pkg
	_pkg, err := pkg.Load(p.Name)
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

	__pkg, err := pkg.Load(p.Name)
	if err != nil {
		return err
	}
	p.Size = __pkg.GetSize()
	p.Dependencies, err = __pkg.GetDependencies()
	if err != nil {
		return err
	}

	return nil
}

func setupPackage(pathToFile string) error {
	return pkg.Extract(pathToFile)
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
		s = s[startQuote+1 : endQuote]
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
