package whl

import (
	"github.com/fextpkg/cli/fext/cfg"
	"github.com/fextpkg/cli/fext/utils"

	"errors"
	"os"
	"strings"
)

const (
	FORMAT_WHEEL = "dist-info"
)

type Package struct {
	Name    string
	metaDir string
	Data    *map[string]string
	Format  string
}

func LoadPackage(name string) (*Package, error) {
	dir, err := findOptimalPackageMetaDir(name)
	if err != nil {
		return nil, err
	} else if dir == "" {
		return nil, errors.New("Package info not found")
	}

	p := Package{Name: name, metaDir: dir}
	s := strings.Split(dir, ".")

	p.Format = s[len(s) - 1]

	return &p, nil
}

// load data from wheel/egg and convert it to single format
func (p *Package) LoadMetaData(libDir string) error {
	var loadFunc func(string) (*map[string]string, error)
	if p.Format == FORMAT_WHEEL {
		loadFunc = loadMeta
	}

	data, err := loadFunc(p.metaDir)
	if err != nil {
		return err
	}
	p.Data = data

	return nil
}

func (p *Package) Uninstall() error {
	dirs := utils.GetAllPackageDirs(p.Name)

	if len(dirs) == 0 {
		return errors.New("Package not installed")
	}

	for _, dir := range dirs {
		err := os.RemoveAll(cfg.PathToLib + dir)
		if err != nil {
			return err
		}
	}

	return nil
}

// Parse dependencies of wheel metadata. Returns error if package have
// unsupported format or another parse error
func (p *Package) GetDependencies() ([]string, error) {
	if p.Format == FORMAT_WHEEL {
		rawDependencies, _, err := loadRawDependenciesAndExtra(p.metaDir)
		if err != nil {
			return nil, err
		}
		return parseDependencies(p.metaDir, rawDependencies)
	} else {
		return nil, errors.New("Unsupported format: " + p.Format)
	}
}

// Get extra packages. Returns error if package have unsupported format
// or another parse error
func (p *Package) GetExtraPackages(names []string) ([]string, error) {
	if p.Format == FORMAT_WHEEL {
		// TODO : do something with errors when extra name not found
		_, rawExtra, err := loadRawDependenciesAndExtra(p.metaDir)
		if err != nil {
			return nil, err
		}

		var extra []string
		for _, name := range names {
			packages, err := parseExtra(name, p.metaDir, rawExtra)
			if err != nil {
				return nil, err
			}
			extra = append(extra, packages...)
		}
		return extra, nil
	} else {
		return nil, errors.New("Unsupported format: " + p.Format)
	}
}

// Calculate all size of files in directory with source code. Returns size in bytes
func (p *Package) GetSize() int64 {
	return utils.GetDirSize(strings.SplitN(p.metaDir, "-", 2)[0])
}
