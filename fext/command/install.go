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
	progressBar *ui.ProgressBar
	optSingle   bool // without dependencies
	optSilent   bool
)

var (
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
		// TODO replace with &
	} else if startQuote != 0 || endQuote != 0 {
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
			return nil, errors.New("extra not found")
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
	pkgName, op := expression.ParseExpression(pkgName)
	web := web.NewRequest(pkgName, op)

	progressBar.UpdateStatus("Scanning", pkgName)
	//progressBar.UpdateStatus("Scanning", pkgName)
	version, link, err := web.GetPackageData()
	if err != nil {
		return err
	}

	p, err := pkg.Load(pkgName)
	if err == nil {
		if version == p.Version {
			return packageAlreadyInstalled
		} else {
			p.Uninstall()
		}
	}

	progressBar.UpdateStatus("Downloading", pkgName)
	filePath, err := web.DownloadPackage(link)
	if err != nil {
		return err
	}

	progressBar.UpdateStatus("Extracting", pkgName)
	if err = io.ExtractPackage(filePath); err != nil {
		return err
	}
	os.RemoveAll(filePath) // remove tmp file, that was downloaded
	if !silent {
		progressBar.Println(fmt.Sprintf("+ %s (%s)", pkgName, version))
	}

	p, err = pkg.Load(pkgName)
	if err != nil {
		return err
	}

	if !optSingle {
		progressBar.UpdateStatus("Installing dependencies")
		for _, dep := range p.Dependencies {
			err = install(fmt.Sprint(dep.Name, dep.Conditions), true)
			if err != nil && err != packageAlreadyInstalled {
				progressBar.Println(fmt.Sprintf("- %s (%s) (%v)", dep.Name, pkgName, err))
			}
		}

	}

	return nil
}

func Install(packages []string) {
	for _, f := range config.Flags {
		switch f {
		case "h", "help":
			ui.PrintHelpInstall()
			return
		case "S", "single":
			optSingle = true
		case "s", "silent":
			optSilent = true
		default:
			ui.PrintUnknownOption(f, ui.PrintHelpInstall)
			return
		}
	}

	progressBar = ui.CreateProgressBar()
	progressBar.UpdateStatus("Parsing packages")
	progressBar.Start()

	for _, pkgName := range packages {
		pkgName, extraNames, err := parseExtraNames(pkgName)
		if err != nil {
			progressBar.Println(fmt.Sprintf("- %s extras (%v)", pkgName, err))
		} else if len(extraNames) > 0 {
			extraPackages, err := getExtraPackages(pkgName, extraNames)
			if err != nil {
				progressBar.Println(fmt.Sprintf("- %s extras (%v)", pkgName, err))
			}

			for _, ePkgName := range extraPackages {
				err = install(ePkgName, optSilent)
				if err != nil {
					progressBar.Println(fmt.Sprintf("- %s (%s) (%v)", ePkgName, pkgName, err))
				}
			}
		} else if err = install(pkgName, optSilent); err != nil {
			progressBar.Println(fmt.Sprintf("- %s (%v)", pkgName, err))
		}
	}
	progressBar.Finish("Finished")
}
