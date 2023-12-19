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
	local []*Query    // Installed packages
	queue chan *Query // Prepared package queries

	opt *Options
}

// supply separates extra dependencies from packages and adds all
// dependent packages to the queue
func (i *Installer) supply(queries []*Query) error {
	var q *Query
	for len(queries) > 0 {
		q = queries[0]
		queries = queries[1:]

		pkgName, extra, err := expression.ParseExtraNames(q.pkgName)
		if err != nil {
			return err
		} else if len(extra) > 0 {
			extraPackages, err := getPackageExtras(pkgName, extra)
			if err != nil {
				if errors.Is(err, ferror.PackageDirectoryMissing) {
					// when trying to install extra dependencies of a package
					// that not installed, install the package and add this query
					// to the end
					// FIXME: separate name and conditions and requeue
					deps, err := i.install(newRawQuery(pkgName))
					if err != nil {
						return err
					}

					queries = append(queries, extrasToQuery(deps)...)
					queries = append(queries, q)
					continue
				}
				return err
			}

			// append extra packages
			queries = append(queries, extraPackages...)
		} else {
			i.queue <- q
		}
	}

	return nil
}

// install installs a single package. Returns its dependencies or an error in
// case of failure
func (i *Installer) install(query *Query) ([]pkg.Dependency, error) {
	req := web.NewRequest(query.pkgName, query.conditions)

	version, link, err := req.GetPackageData()
	if err != nil {
		return nil, err
	}

	// check if package already installed
	p, err := pkg.Load(query.pkgName)
	if err == nil {
		if version == p.Version {
			return nil, ferror.PackageAlreadyInstalled
		} else {
			if err = p.Uninstall(); err != nil {
				return nil, err
			}
		}
	}

	filePath, err := req.DownloadPackage(link)
	if err != nil {
		return nil, err
	}

	if err = io.ExtractPackage(filePath); err != nil {
		return nil, err
	}
	// remove downloaded file
	if err = os.RemoveAll(filePath); err != nil {
		return nil, err
	}

	// check the package installed correctly
	p, err = pkg.Load(query.pkgName)
	if err != nil {
		return nil, err
	}

	// Make a note that fext installed this package
	err = io.CreateInstallerFile(p.GetMetaDirectoryPath())
	if err != nil {
		return nil, err
	}

	return p.GetDependencies(), nil
}

// process pops the package from queue and installs it. Parses dependencies
// of installed package and append them to queue. Prints the final result
func (i *Installer) process() {
	for len(i.queue) > 0 {
		q, open := <-i.queue
		if !open {
			break
		}

		dependencies, err := i.install(q)
		if err != nil {
			ui.PrintfMinus("%s (%v)\n", q.pkgName, err)
			continue
		}

		if !i.opt.QuietMode {
			ui.PrintlnPlus(q.pkgName)
		}

		if !i.opt.NoDependencies {
			err = i.supply(extrasToQuery(dependencies))
			if err != nil {
				ui.PrintfError("%s deps (%s)\n", q.pkgName, err)
			}
		}
	}
}

// InitializePackages converts the list of packages into a query list, parses
// extra dependencies and adds the prepared packages to the queue
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
		e, err := p.GetExtraDependencies(extraName)
		if err != nil {
			return nil, err
		} else if len(e) == 0 {
			return nil, &ferror.MissingExtra{Name: extraName}
		}

		for _, dep := range e {
			queries = append(queries, newQuery(dep.PackageName, dep.Conditions))
		}
	}

	return queries, nil
}
