package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
)

var (
	MissingPackage    = errors.New("package not found")
	UnsupportedFormat = errors.New("unsupported format")
)

type Package struct {
	metaDir string
	Name    string
	Version string

	Dependencies []extra
	Extra        map[string][]extra
}

func Load(pkgName string) (*Package, error) {
	dirName, err := getPackageMetaDir(pkgName)
	if err != nil {
		return nil, err
	}

	p := Package{
		Name:         pkgName,
		metaDir:      dirName,
		Dependencies: []extra{},
		Extra:        map[string][]extra{},
	}
	if err = p.parseMetaData(); err != nil {
		return nil, err
	}

	return &p, nil
}

func LoadFromMetaDir(metaDir string) (*Package, error) {
	p := Package{
		metaDir:      metaDir,
		Dependencies: []extra{},
		Extra:        map[string][]extra{},
	}
	if err := p.parseMetaData(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Package) parseMetaData() error {
	data, err := os.ReadFile(getAbsolutePath(p.metaDir, "METADATA"))
	if err != nil {
		return err
	}

	var extraName string
	for _, s := range strings.Split(strings.SplitN(string(data), "\n\n", 2)[0], "\n") {
		// FIXME: this a temporary solution that will be rewritten in the future
		if s != "" && (s[0] == 'R' || s[0] == 'P' || s[0] == 'V' || s[0] == 'N') {
			field := strings.SplitN(s, ": ", 2)
			if field[0] == "Requires-Dist" {
				e := extra{}
				value := strings.Split(field[1], ";") // [name_and_conditions, markers]
				if len(value) == 2 {
					e.markers = value[1]
				}
				value = strings.Split(value[0], " ") // [name, conditions]
				if len(value) > 1 {
					e.Conditions = value[1]
				}
				e.Name = value[0]

				if extraName != "" {
					if _e, found := p.Extra[extraName]; found {
						_e = append(_e, e)
						p.Extra[extraName] = _e
					} else {
						p.Extra[extraName] = []extra{e}
					}
				} else {
					p.Dependencies = append(p.Dependencies, e)
				}
			} else if field[0] == "Provides-Extra" {
				extraName = field[1]
			} else if field[0] == "Version" {
				p.Version = strings.Replace(field[1], "\r", "", 1)
			} else if field[0] == "Name" {
				p.Name = strings.Replace(field[1], "\r", "", 1)
			}
		}
	}
	return nil
}

// getTopLevel returns the packages names from the file "top_level.txt". This
// file stores a list of packages that python can work with
func (p *Package) getTopLevel() ([]string, error) {
	var packages []string
	data, err := os.ReadFile(getAbsolutePath(p.metaDir, "top_level.txt"))
	if err != nil {
		return nil, err
	}

	packages = strings.Split(strings.TrimSpace(string(data)), "\n")
	return packages, nil
}

// getSourceDirs returns a list of directories with source code
func (p *Package) getSourceDirs() ([]string, error) {
	packages, err := p.getTopLevel()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		// additionally, we check the presence of the source code folder, since some
		// generators do not add the top_level.txt file
		if _, err = os.Stat(getAbsolutePath(p.Name)); err == nil {
			packages = []string{p.Name}
		}
	}
	return packages, nil
}

// Uninstall deletes all folders and files associated with this package
func (p *Package) Uninstall() error {
	packages, err := p.getSourceDirs()
	if err != nil {
		return err
	}

	removeDir := func(dirName string) error { return os.RemoveAll(getAbsolutePath(dirName)) }
	if len(packages) == 0 { // this is not a package but a module
		if err = os.Remove(getAbsolutePath(fmt.Sprintf("%s.py", formatName(p.Name)))); err != nil {
			return err
		}
	} else {
		for _, pkgName := range packages {
			if err = removeDir(pkgName); err != nil {
				return err
			}
		}
	}

	if removeDir(p.metaDir) != nil {
		return err
	}

	return nil
}

// GetSize calculate all size of files in directory with source code. Returns size in bytes
func (p *Package) GetSize() (int64, error) {
	packages, err := p.getSourceDirs()
	if err != nil {
		return 0, err
	} else if len(packages) == 0 { // this is not a package but a module
		f, err := os.Stat(getAbsolutePath(formatName(p.Name) + ".py"))
		if err != nil {
			return 0, err
		} else {
			return f.Size(), nil
		}
	}

	var size int64
	for _, pkgName := range packages {
		err = filepath.Walk(getAbsolutePath(pkgName), func(_ string, info os.FileInfo, _ error) error {
			if info != nil && !info.IsDir() {
				size += info.Size()
			}
			return nil
		})
	}

	return size, nil
}

// extra is used simultaneously for dependencies and extra packages
type extra struct {
	Name       string
	Conditions string
	markers    string
}

// CheckMarkers checks the possibility of installation according to the
// specified markers. Returns an error if parsing failed
func (e *extra) CheckMarkers() (bool, error) {
	// TODO move marker replaces from markers module to this func
	return expression.CompareExpression(e.markers)
}

// formatName formats the directory name to a single view
func formatName(dirName string) string {
	return strings.ToLower(strings.ReplaceAll(dirName, "-", "_"))
}

// parseFormat parse the directory name and returns its format.
// Example: "requests-2.26.0.dist-info" => "dist-info"
func parseFormat(dirName string) string {
	s := strings.Split(dirName, ".")
	return s[len(s)-1]
}

// clearVersion removes extra characters from the version.
// Example: "2.26.0.dist" => "2.26.0"
func clearVersion(version string) string {
	return strings.Replace(version, ".dist", "", 1)
}

// getAbsolutePath returns absolute path to the file in directory with packages
func getAbsolutePath(elem ...string) string {
	return filepath.Clean(config.PythonLibPath) + string(os.PathSeparator) + filepath.Join(elem...)
}

// Parse directory by format "%pkgName%-%version%.%format%" and returns it
func parseDirectoryName(dirName string) (string, string, string) {
	// [name, version, format]
	meta := strings.SplitN(dirName, "-", 3)

	// avoid errors
	if len(meta) >= 2 {
		return meta[0], clearVersion(meta[1]), parseFormat(dirName)
	} else {
		return meta[0], "", ""
	}
}

// getPackageMetaDir searches for a folder with wheel format from the specified
// pkgName. Returns the original directory name. Returns an error if the found
// format is not supported or the package is missing
func getPackageMetaDir(pkgName string) (string, error) {
	dirInfo, err := os.ReadDir(config.PythonLibPath)
	if err != nil {
		return "", err
	}
	pkgName = formatName(pkgName)

	for _, dir := range dirInfo {
		curPkgName, _, format := parseDirectoryName(dir.Name())
		if formatName(curPkgName) == pkgName && format == "dist-info" {
			return dir.Name(), nil
		}
	}

	return "", MissingPackage
}
