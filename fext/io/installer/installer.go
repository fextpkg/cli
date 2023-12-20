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
	local []*Query
	// Prepared installation queries for the packages, including names,
	// conditional operators, and versions
	queue chan *Query

	opt *Options
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
			}
			// Replacing the package name with a clean one,
			// excluding any extra names
			q.pkgName = pkgName
			// Adding extra dependencies to the queue for further processing,
			// as they may also have additional extra dependencies within them
			queries = append(queries, extraDeps...)
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
		if version == p.Version {
			// The required version of the package is already installed,
			// so we don't need to download it again
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
		if err != nil {
			// Installation of the package failed. Displaying the error message
			// regardless of the QuietMode setting
			ui.PrintfMinus("%s (%v)\n", q.pkgName, err)
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
		} else if len(extraDeps) == 0 {
			return nil, &ferror.MissingExtra{Name: extraName}
		}

		for _, dep := range extraDeps {
			queries = append(queries, newQuery(dep.PackageName, dep.Conditions))
		}
	}

	return queries, nil
}
