package command

import (
	"strings"

	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

type CheckPackageHealth struct {
	metaDirectories []string
}

func InitCheckPackageHealth() *CheckPackageHealth {
	return &CheckPackageHealth{}
}

// scanMissingDependencies scans the dependencies of a package (including
// extras), and if any of them fail to load, they are added to the list of
// missing dependencies.
// Returns a list of package names that could not be loaded, or error if any
// package was failed to load.
func (cmd *CheckPackageHealth) scanMissingDependencies(p *pkg.Package) ([]string, error) {
	var missingDependencies []string
	packages, err := p.GetDependencies()
	if err != nil {
		return nil, err
	}

	for len(packages) > 0 {
		dep := packages[0]
		packages = packages[1:]

		pkgName, extraNames, err := expression.ParseExtraNames(dep.PackageName)
		if err != nil {
			return nil, err
		}

		depPackage, err := pkg.Load(pkgName)
		if err != nil {
			missingDependencies = append(missingDependencies, pkgName)
		} else if len(extraNames) > 0 {
			for _, extra := range extraNames {
				extraPackages, err := depPackage.GetExtraDependencies(extra)
				if err != nil {
					return nil, err
				}
				// Since extra packages can contain additional extra packages,
				// we will implement their verification through a list to avoid
				// recursion
				packages = append(packages, extraPackages...)
			}
		}
	}

	return missingDependencies, nil
}

// scanMismatchingDependencies scans the dependencies of a package and only
// checks their version compatibility if the package and its dependencies are
// loaded correctly.
// Returns a list of package names where incompatibilities are found,
// or error if the version comparison resulted in an error.
func (cmd *CheckPackageHealth) scanMismatchingDependencies(p *pkg.Package) ([]string, error) {
	var deps []string

	dependencies, err := p.GetDependencies()
	if err != nil {
		return nil, err
	}

	for _, dep := range dependencies {
		p, err := pkg.Load(dep.PackageName)
		if err != nil {
			// Skipping the package, as method is not responsible for missing
			// dependencies
			continue
		}

		result, err := expression.CompareConditions(p.Version, dep.Conditions)
		if err != nil {
			// It is better to return an error to explicitly indicate the
			// issues in the system rather than ignoring it
			return nil, err
		} else if !result {
			deps = append(deps, dep.PackageName)
		}
	}

	return deps, nil
}

// checkPackageDependencies checks a package for installation errors. If the
// package has incompatible versions of dependencies, or they are missing
// altogether, an error will be displayed. If the package fails to load,
// an error will also be displayed.
// Returns the total number of missing and incompatible dependencies.
func (cmd *CheckPackageHealth) checkPackageDependencies(metaDir string) (int, error) {
	p, err := pkg.LoadFromMetaDir(metaDir)
	if err != nil {
		return 1, err
	}

	missingDeps, err := cmd.scanMissingDependencies(p)
	if err != nil {
		return 1, err
	} else if len(missingDeps) > 0 {
		ui.PrintfError(
			"check %s: missing dependencies: %s\n",
			p.Name,
			strings.Join(missingDeps, ", "),
		)
	}

	mismatchingDeps, err := cmd.scanMismatchingDependencies(p)
	if err != nil {
		return 1, err
	} else if len(mismatchingDeps) > 0 {
		ui.PrintfError(
			"check %s: mismatching dependencies versions: %s\n",
			p.Name,
			strings.Join(mismatchingDeps, ", "),
		)
	}

	return len(missingDeps) + len(mismatchingDeps), nil
}

// DetectFlags does nothing and is a stub to maintain a single interface of
// interaction.
func (cmd *CheckPackageHealth) DetectFlags() error {
	return nil
}

// Execute has iterates through all packages installed in the system and check
// if everything is fine with them. If any issues are found with a package
// (incompatibility with dependencies, missing packages, or failed to load),
// an error message will be displayed. Otherwise, if everything is fine,
// an "ok" message will be displayed.
func (cmd *CheckPackageHealth) Execute() {
	var brokenPackages int
	var err error

	cmd.metaDirectories, err = io.GetMetaDirectories()
	if err != nil {
		ui.Fatal("Unable to scan meta directories: " + err.Error())
	}

	for _, dirName := range cmd.metaDirectories {
		brokenCount, err := cmd.checkPackageDependencies(dirName)
		if err != nil {
			ui.PrintfError("check %s: error: %v\n", dirName, err.Error())
		}

		brokenPackages += brokenCount
	}

	if brokenPackages == 0 {
		ui.PrintfOK("Everything is ok\n")
	}
}
