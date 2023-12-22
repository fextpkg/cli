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

func InitShowPackageInfo(args []string) *ShowPackageInfo {
	return &ShowPackageInfo{
		packageNames: args,
	}
}

// DetectFlags does nothing and is a stub to maintain a single interface of
// interaction.
func (cmd *ShowPackageInfo) DetectFlags() error {
	return nil
}

// Execute prints general information about the first package
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

// prettifyData loads a package and returns information about it,
// formatted nicely and in user-friendly manner.
// It returns an error if it fails to process the data for the given pkgName.
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

// checkPackageInstalled checks if a package is installed by loading and
// processing it.
func checkPackageInstalled(pkgName string) bool {
	_, err := pkg.Load(pkgName)
	if err != nil {
		return false
	}
	return true
}

// prettifyDependencies formats the list of dependencies in a visually
// appealing and user-friendly format. Package names will be colored green if
// the dependency is installed and everything is fine. They will be a colored
// red if the dependency is missing or if there was an error during the loading
// process. It returns a dash (-) if no dependencies are found.
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

// prettifyExtraDependencies formats the list of extra dependency names into a
// visually appealing and user-friendly format. Names will be colored to green
// if all packages from the extra dependency list are installed correctly.
// They will be colored to red if there is an error during parsing.
// They will have no color if they are not installed in the system.
// It returns a dash (-) if no extras are found.
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
