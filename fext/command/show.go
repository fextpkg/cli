package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

type ShowPackageInfo struct {
	packageNames []string
}

func (cmd *ShowPackageInfo) DetectFlags() error {
	return nil
}

func (cmd *ShowPackageInfo) Execute() {
	if len(cmd.packageNames) == 0 {
		ui.PrintlnError("Unable to get package data: no package were passed")
		return
	}

	data, err := prettifyData(cmd.packageNames[0])
	if err != nil {
		ui.PrintlnError("Unable to get package data: " + err.Error())
		return
	}

	fmt.Println(data)
}

func InitShowPackageInfo(args []string) *ShowPackageInfo {
	return &ShowPackageInfo{
		packageNames: args,
	}
}

func prettifyData(pkgName string) (string, error) {
	p, err := pkg.Load(pkgName)
	if err != nil {
		return "", err
	}

	size, err := p.GetSize()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Name: %s\nVersion: %s\nSize: %s MB\nDependencies: %s\nExtra Dependencies: %s",
		ui.BoldString(p.Name),
		ui.BoldString(p.Version),
		ui.BoldString(strconv.FormatFloat(float64(size/1024)/1024, 'f', 2, 32)),
		prettifyDependencies(p.GetDependencies()),
		prettifyExtraDependencies(p),
	), nil
}

func checkPackageInstalled(pkgName string) bool {
	_, err := pkg.Load(pkgName)
	if err != nil {
		return false
	}
	return true
}

func prettifyDependencies(deps []pkg.Dependency) string {
	var text strings.Builder

	for _, dep := range deps {
		_, err := pkg.Load(dep.PackageName)
		if err != nil {
			text.WriteString(ui.RedString(dep.PackageName))
		} else {
			text.WriteString(ui.GreenString(dep.PackageName))
		}
		text.WriteString(", ")
	}

	if text.Len() == 0 {
		return "-"
	}

	return text.String()[:text.Len()-2]
}

func prettifyExtraDependencies(p *pkg.Package) string {
	var text strings.Builder

	for _, extraName := range p.Extras {
		deps, err := p.GetExtraDependencies(extraName)
		if err != nil {
			text.WriteString(ui.RedString(extraName))
		} else {
			var installed = true
			for _, dep := range deps {
				if !checkPackageInstalled(dep.PackageName) {
					installed = false
				}
			}

			if installed {
				text.WriteString(ui.GreenString(extraName))
			} else {
				text.WriteString(extraName)
			}
		}

		text.WriteString(", ")
	}

	if text.Len() == 0 {
		return "-"
	}

	return text.String()[:text.Len()-2]
}
