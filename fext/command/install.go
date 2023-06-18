package command

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/io/web"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

var (
	optNoDependencies bool // without dependencies
	optSilent         bool // output only error messages

	packageAlreadyInstalled = errors.New("package already installed")
)

// Parse extra names. Returns pkgName, extraNames and error if syntax is invalid
func parseExtraNames(s string) (string, []string, error) {
	var startQuote, endQuote int
	for i, v := range s {
		// find first quote
		if v == '[' && startQuote == 0 {
			startQuote = i
			// find last quote
		} else if v == ']' {
			endQuote = i
		}
	}

	if startQuote != 0 && endQuote != 0 {
		originalName := s[:startQuote]
		s = s[startQuote+1 : endQuote]
		if strings.ContainsAny(s, "[]") {
			return originalName, nil, errors.New("syntax error")
		}

		var extraNames []string
		for _, name := range strings.Split(s, ",") {
			name = strings.ReplaceAll(name, " ", "")
			extraNames = append(extraNames, name)
		}

		return originalName, extraNames, nil
	} else if startQuote|endQuote != 0 {
		return s, nil, errors.New("syntax error")
	}

	return s, nil, nil
}

func getExtraPackages(pkgName string, extraNames []string) ([]string, error) {
	var packages []string
	p, err := pkg.Load(pkgName)
	if err != nil {
		return nil, err
	}
	for _, eName := range extraNames {
		e, ok := p.Extra[eName]
		if !ok {
			return nil, errors.New("extra not found: " + eName)
		}
		for _, extra := range e {
			passMarkers, err := extra.CheckMarkers()
			if err != nil {
				return nil, err
			} else if !passMarkers {
				continue
			}
			packages = append(packages, fmt.Sprint(extra.Name, extra.Conditions))
		}
	}
	return packages, nil
}

func install(pkgName string, silent bool) error {
	pkgName, conditions := expression.ParseConditions(pkgName)
	web := web.NewRequest(pkgName, conditions)

	version, link, err := web.GetPackageData()
	if err != nil {
		return err
	}

	p, err := pkg.Load(pkgName)
	if err == nil {
		if version == p.Version {
			return packageAlreadyInstalled
		} else {
			if err = p.Uninstall(); err != nil {
				return err
			}
		}
	}

	filePath, err := web.DownloadPackage(link)
	if err != nil {
		return err
	}

	if err = io.ExtractPackage(filePath); err != nil {
		return err
	}
	if err = os.RemoveAll(filePath); err != nil { // remove tmp file, that was downloaded
		return err
	}
	if !silent {
		ui.PrintfPlus("%s (%s)\n", pkgName, version)
	}

	p, err = pkg.Load(pkgName)
	if err != nil {
		return err
	}

	if !optNoDependencies {
		for _, dep := range p.Dependencies {
			err = install(fmt.Sprint(dep.Name, dep.Conditions), true)
			if err != nil && err != packageAlreadyInstalled {
				ui.PrintfMinus("%s (%s) (%v)\n", dep.Name, pkgName, err)
			}
		}

	}

	return nil
}

func getPackagesFromFiles(files []string) ([]string, error) {
	var packages []string
	for _, fileName := range files {
		p, err := io.ReadLines(fileName)
		if err != nil {
			return nil, err
		}
		packages = append(packages, p...)
	}
	return packages, nil
}

func Install(packages []string) {
	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			ui.PrintHelpInstall()
			return
		case "n", "no-dependencies":
			optNoDependencies = true
		case "s", "silent":
			optSilent = true
		case "r", "requirements":
			var err error
			packages, err = getPackagesFromFiles(packages)
			if err != nil {
				ui.PrintlnError("Failed to read file:", err.Error())
			}
		default:
			ui.PrintUnknownOption(f, ui.PrintHelpInstall)
			return
		}
	}

	for _, pkgName := range packages {
		pkgName, extraNames, err := parseExtraNames(pkgName)
		if err != nil {
			ui.PrintfMinus("%s extras (%v)\n", pkgName, err)
		} else if len(extraNames) > 0 {
			extraPackages, err := getExtraPackages(pkgName, extraNames)
			if err != nil {
				ui.PrintfMinus("%s extras (%v)\n", pkgName, err)
			}

			for _, ePkgName := range extraPackages {
				err = install(ePkgName, optSilent)
				if err != nil {
					ui.PrintfMinus("%s (%s) (%v)\n", ePkgName, pkgName, err)
				}
			}
		} else if err = install(pkgName, optSilent); err != nil {
			ui.PrintfMinus("%s (%v)\n", pkgName, err)
		}
	}
}
