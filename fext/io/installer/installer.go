package installer

import (
	"errors"
	"os"

	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/io/web"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

type Options struct {
	// Do not install package dependencies
	NoDependencies bool

	// Output only error messages
	QuietMode bool
}

// DefaultOptions returns an Options struct with default parameters
func DefaultOptions() *Options {
	return &Options{
		NoDependencies: false,
		QuietMode:      false,
	}
}

type Installer struct {
	// Installed packages
	local map[string]*Query
	// Prepared installation queries for the packages, including names,
	// conditional operators, and versions
	queue chan *Query

	opt *Options
}

// updateLocal updates or adds conditional operators to maintain compatibility
// with the other packages to be installed.
func (i *Installer) updateLocal(newQuery *Query) {
	query, exist := i.local[newQuery.pkgName]
	if exist {
		query.conditions = append(query.conditions, newQuery.conditions...)
	} else {
		i.local[newQuery.pkgName] = newQuery
	}
}

// checkCompatibility checks the compatibility of the package version with
// other packages that have been installed in the current session. It compares
// the existing conditions with the ones provided.
// It returns a boolean value indicating the compatibility or false if the
// package was not found in local. If an error occurs during the comparison of
// operators, it throws an error.
func (i *Installer) checkCompatibility(pkgName, version string, cond []expression.Condition) (bool, error) {
	query, exist := i.local[pkgName]
	if exist {
		return expression.CompareConditions(version, append(query.conditions, cond...))
	}

	return false, nil
}

// isInstalled checks if the package has been installed within the
// current session.
func (i *Installer) isInstalled(pkgName string) bool {
	_, exist := i.local[pkgName]
	return exist
}

// Extracts extra dependencies from the packages, creates a new query object
// with them, and adds it to the queue. If an extra package name is provided
// but not locally installed, the query object is copied to the "extraNames"
// attribute, and the original query is cleared of extra names.
// It returns an error if the extra packages are not found or if there is any
// other issue related to processing the package metadata.
func (i *Installer) supply(queries []*Query) error {
	var q *Query
	for len(queries) > 0 {
		q = queries[0]
		queries = queries[1:]

		pkgName, extra, err := expression.ParseExtraNames(q.pkgName)
		if err != nil {
			return err
		} else if len(extra) > 0 {
			extraDeps, err := getPackageExtras(pkgName, extra)
			if err != nil {
				if !errors.Is(err, ferror.PackageDirectoryMissing) {
					return err
				}
				// If a locally installed package is not found, but we are
				// attempting to install its extra dependencies, an error would
				// normally occur. To prevent this error and maintain backward
				// compatibility, we employ a workaround by duplicating the
				// query object without the extra names, and instead add the
				// current query to the "extraNames" attribute. Once the
				// installation is complete, we can retrieve and process the
				// extra dependencies accordingly
				q.extraNames = copyQuery(q)
			} else {
				// Adding extra dependencies to the queue for further processing,
				// as they may also have additional extra dependencies within them
				queries = append(queries, extraDeps...)
				if len(q.conditions) == 0 {
					// If no installation conditions are specified,
					// it indicates that the package is already installed.
					// Hence, there is no need for additional installation
					continue
				}
			}
			// Replacing the package name with a clean one,
			// excluding any extra names
			q.pkgName = pkgName
		}
		i.queue <- q
	}

	return nil
}

// Fetches the available versions for installation, selects the suitable
// version based on the provided query attributes, downloads and unpacks the
// package into the config.PythonLibPath.
// It returns the package dependencies or an error if any occurs.
func (i *Installer) install(query *Query) ([]pkg.Dependency, error) {
	// Creating a new request
	req := web.NewRequest(query.pkgName, query.conditions)

	// Retrieving the necessary version based on the provided parameters
	version, link, err := req.GetPackageData()
	if err != nil {
		return nil, err
	}

	// First, check if the package is installed locally
	p, err := pkg.Load(query.pkgName)
	if err == nil {
		compatible, err := i.checkCompatibility(query.pkgName, p.Version, query.conditions)
		if err != nil {
			// An error occurred while comparing operators
			return nil, err
		} else if compatible {
			// The package is already installed and compatible with other
			// packages that rely on it. Therefore, it is not necessary to
			// reinstall it
			return nil, ferror.PackageInLocalList
		} else if version == p.Version && !i.isInstalled(query.pkgName) {
			// This is the initial request for package installation.
			// We need to provide a meaningful error to explain what occurred
			return nil, ferror.PackageAlreadyInstalled
		} else {
			// The package is installed, but the version is not suitable.
			// Remove the package and proceed with installing the required version
			if err = p.Uninstall(); err != nil {
				return nil, err
			}
		}
	}

	// Commencing package download
	filePath, err := req.DownloadPackage(link)
	if err != nil {
		return nil, err
	}

	// Unpacking the installed file
	if err = io.ExtractPackage(filePath); err != nil {
		return nil, err
	}
	// Remove downloaded file
	if err = os.RemoveAll(filePath); err != nil {
		return nil, err
	}

	// Finally, we ensure that the package is installed correctly
	p, err = pkg.Load(query.pkgName)
	if err != nil {
		return nil, err
	}

	// Make a note that fext installed this package
	if err = io.CreateInstallerFile(p.GetMetaDirectoryPath()); err != nil {
		return nil, err
	}

	return p.GetDependencies(), nil
}

// Retrieves a package from the queue and starts its installation along
// with its dependencies. It also installs extra dependencies if they were
// provided. The output is displayed in stdout. If Options.QuietMode if set to
// true, success messages will not be displayed.
func (i *Installer) process() {
	for len(i.queue) > 0 {
		q, open := <-i.queue
		if !open {
			// Handling exceptional cases when we abruptly close the channel
			break
		}

		dependencies, err := i.install(q)
		// Update conditions even if an error occurs, because there is hope for
		// a subsequent installation request.
		i.updateLocal(q)
		if err != nil {
			if !errors.Is(err, ferror.PackageInLocalList) && !(errors.Is(err, ferror.PackageAlreadyInstalled) && q.isDependency) {
				// The condition is passed only if the package has not been
				// installed before within the current session. Or if another
				// error is received and the package is not a dependency.
				// If the package is a dependency and is already installed,
				// there is no point in displaying an error message. Otherwise,
				// there is no need to conceal the error
				ui.PrintfMinus("%s (%v)\n", q.pkgName, err)
			}
			continue
		}

		if !i.opt.QuietMode {
			// Displaying a success message only if the quiet mode is not enabled
			ui.PrintlnPlus(q.pkgName)
		}

		if !i.opt.NoDependencies {
			// Installing the acquired package dependencies during installation
			err = i.supply(extraDependenciesToQuery(dependencies))
			if err != nil {
				ui.PrintfMinus("%s deps (%s)\n", q.pkgName, err)
			}
		}

		if q.extraNames != nil {
			// Preparing extra dependencies for installation in this run due to
			// their absence in the system during the previous attempt
			err = i.supply([]*Query{q.extraNames})
			if err != nil {
				ui.PrintfMinus("%s extras (%s)\n", q.pkgName, err)
			}
		}
	}
}

// InitializePackages converts package names into a query queue and adds them
// to the queue using the supply method.
// It returns any error returned by the supply method.
func (i *Installer) InitializePackages(packages []string) error {
	var q []*Query
	for _, pkgName := range packages {
		q = append(q, newRawQuery(pkgName))
	}

	return i.supply(q)
}

// Install starts the package installing loop
func (i *Installer) Install() {
	i.process()
	close(i.queue)
}

func NewInstaller(opt *Options) *Installer {
	return &Installer{
		local: map[string]*Query{},
		queue: make(chan *Query, 512),
		opt:   opt,
	}
}

// getPackageExtras gets the extra dependencies of the package and wraps them in
// a query list. Returns ferror.MissingExtra if extra name not found
func getPackageExtras(pkgName string, extraNames []string) ([]*Query, error) {
	var queries []*Query
	p, err := pkg.Load(pkgName)
	if err != nil {
		return nil, err
	}

	for _, extraName := range extraNames {
		extraDeps, err := p.GetExtraDependencies(extraName)
		if err != nil {
			return nil, err
		} else if !p.HasExtraName(extraName) {
			return nil, &ferror.MissingExtra{Name: extraName}
		}

		for _, dep := range extraDeps {
			queries = append(queries, newQuery(dep.PackageName, dep.Conditions, true))
		}
	}

	return queries, nil
}
