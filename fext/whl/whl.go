package whl

import (
	"errors"
	"github.com/Flacy/fext/fext/base_cfg"
	"github.com/Flacy/fext/fext/utils"
	"os"
	"path/filepath"
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

func LoadPackage(name, libDir string) (*Package, error) {
	dir, err := findOptimalPackageMetaDir(name, libDir)
	if err != nil {
		return nil, err
	} else if dir == "" {
		return nil, errors.New("Package info not found")
	}

	p := Package{Name: name, metaDir: libDir + dir}
	s := strings.Split(dir, ".")

	p.Format = s[len(s) - 1]

	return &p, nil
}

// load data from wheel/egg and convert it to single format
func (p *Package) LoadMetaData(libDir string) error {
	var loadFunc func(string) (*map[string]string, error)
	if p.Format == FORMAT_WHEEL {
		loadFunc = loadWheelMeta
	}

	data, err := loadFunc(p.metaDir)
	if err != nil {
		return err
	}
	p.Data = data

	return nil
}

func (p *Package) LoadDependencies() []string {
	var loadFunc func(string) []string
	if p.Format == FORMAT_WHEEL {
		loadFunc = loadWheelDependencies
	} else if p.Format == FORMAT_EGG {
		loadFunc = loadEggDependencies
	}

	return loadFunc(p.metaDir)
}

func (p *Package) Uninstall(libDir string) error {
	dirs := utils.GetAllPackageDirs(p.Name, libDir)

	if len(dirs) == 0 {
		return errors.New("Package not installed")
	}

	for _, dir := range dirs {
		err := os.RemoveAll(libDir + dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Package) GetSize() int64 {
	// cut part with meta and append original name
	return utils.GetDirSize(filepath.Dir(p.metaDir) + base_cfg.PATH_SEPARATOR + p.Name)
}
