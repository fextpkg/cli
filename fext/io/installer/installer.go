package installer

import (
	"os"

	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/io"
	"github.com/fextpkg/cli/fext/io/web"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

type Installer struct {
	local []*Query    // Installed packages
	queue chan *Query // Prepared package queries
}

// supplyPackages separates extra dependencies from packages and adds all
// prepared packages to the queue
func (i *Installer) supplyPackages(queries []*Query) error {
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
				if err == ferror.PackageDirectoryMissing {
					// when trying to install extra dependencies of a package
					// that not installed, install the package and add this query
					// to the end
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
			for _, extraQuery := range extraPackages {
				queries = append(queries, extraQuery)
			}
		} else {
			i.queue <- q
		}
	}

	return nil
}

// install installs a single package. Returns its dependencies or an error in
// case of failure
func (i *Installer) install(query *Query) ([]pkg.Extra, error) {
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
	// remove installed file
	if err = os.RemoveAll(filePath); err != nil {
		return nil, err
	}

	// check the package installed correctly
	p, err = pkg.Load(query.pkgName)
	if err != nil {
		return nil, err
	}

	return p.Dependencies, nil
}

// process pops the package from queue and installs it. Parses dependencies
// and append them to queue
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
		ui.PrintlnPlus(q.pkgName)

		err = i.supplyPackages(extrasToQuery(dependencies))
		if err != nil {
			ui.PrintfError("%s deps (%s)\n", q.pkgName, err)
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

	return i.supplyPackages(q)
}

// Install starts the package installing loop
func (i *Installer) Install() {
	i.process()
	close(i.queue)
}

func NewInstaller() *Installer {
	return &Installer{
		queue: make(chan *Query, 512),
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
		e, ok := p.Extra[extraName]
		if !ok {
			return nil, ferror.NewMissingExtra(extraName)
		}
		for _, extra := range e {
			if extra.Compatible {
				queries = append(queries, newQuery(extra.Name, extra.Conditions))
			}
		}
	}

	return queries, nil
}
