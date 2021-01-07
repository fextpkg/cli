package whl

import (
	"errors"
	"github.com/Flacy/fext/fext/cfg"
	"github.com/Flacy/fext/fext/utils"
	"os"
	"strings"
)

const (
	FORMAT_WHEEL = "dist-info"
	FORMAT_EGG = "egg-info"
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

// Parse dependencies of wheel metadata. Returns error if package have unsupported format
func (p *Package) GetDependencies() ([]string, error) {
	if p.Format == FORMAT_WHEEL {
		return loadDependencies(p.metaDir)
	} else {
		return nil, errors.New("Unsupported format: " + p.Format)
	}
}

// Get extra packages. Returns error if extra name doesn't exists, or another parse error
func (p Package) GetExtra(name string) ([]string, error) {
	// TODO
	return nil, nil
}

// Calculate all size of files in directory with source code. Returns size in bytes
func (p *Package) GetSize() int64 {
	// TODO check workable
	return utils.GetDirSize(strings.SplitN(p.metaDir, "-", 2)[0])
}
