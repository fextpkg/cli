package whl

import (
	"errors"
	"strings"
)

const (
	FORMAT_WHEEL = "dist-info"
	FORMAT_EGG = "egg-info"
)

type Package struct {
	Name string
	dir string
	Data *map[string]string
	Format string
}

func LoadPackage(name, libDir string) (*Package, error) {
	dir, err := findOptimalPackageMetaDir(name, libDir)
	if err != nil {
		return nil, err
	} else if dir == "" {
		return nil, errors.New("Package info not found")
	}

	p := Package{Name: name, dir: libDir + dir}
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

	data, err := loadFunc(p.dir)
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

	return loadFunc(p.dir)
}
