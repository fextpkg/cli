package pkg

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/io"
)

type Package struct {
	metaDir string
	Name    string
	Version string

	Dependencies []Extra
	Extra        map[string][]Extra
}

func Load(pkgName string) (*Package, error) {
	dirName, err := getPackageMetaDir(pkgName)
	if err != nil {
		return nil, err
	}

	p := Package{
		Name:         pkgName,
		metaDir:      dirName,
		Dependencies: []Extra{},
		Extra:        map[string][]Extra{},
	}
	if err = p.parseMetaData(); err != nil {
		return nil, err
	}

	return &p, nil
}

func LoadFromMetaDir(metaDir string) (*Package, error) {
	p := Package{
		metaDir:      metaDir,
		Dependencies: []Extra{},
		Extra:        map[string][]Extra{},
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
				e := Extra{Compatible: true}
				value := strings.Split(field[1], ";") // [name_and_conditions, markers]
				if len(value) == 2 {
					e.Compatible, err = expression.CompareExpression(value[1])
					if err != nil {
						return err
					}
				}
				value = strings.Split(value[0], " ") // [name, conditions]
				if len(value) > 1 {
					_, e.Conditions = expression.ParseConditions(value[1])
				}
				e.Name = value[0]

				if extraName != "" {
					if _e, found := p.Extra[extraName]; found {
						p.Extra[extraName] = append(_e, e)
					} else {
						p.Extra[extraName] = []Extra{e}
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

// getTopLevel scans the "top_level.txt" file, which contains the names of
// packages and modules. If it does not exist, then adds the package name and
// returns it.
func (p *Package) getTopLevel() ([]string, error) {
	files, err := io.ReadLines(getAbsolutePath(p.metaDir, "top_level.txt"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		// add the package name manually, since some generators do not create a
		// "top_level.txt" file
		files = []string{p.Name}
	}
	return files, nil
}

// getSourceFiles returns name of source file belonging to this package by
// converting names of modules and packages
func (p *Package) getSourceFiles() ([]string, error) {
	files, err := p.getTopLevel()
	if err != nil {
		return nil, err
	}
	for i, fileName := range files {
		// since the names of the files are contained without an extension, first we
		// check for a directory with this name, if it is not there, then it is a python
		// file
		if _, err = os.Stat(getAbsolutePath(fileName)); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, err
			}
			files[i] = fileName + ".py"
		}
	}
	return files, nil
}

// getDataDirectory returns the name of the directory with the data files
func (p *Package) getDataDirectory() string {
	return p.metaDir[:len(p.metaDir)-9] + "data"
}

// Uninstall deletes all directories and files belonging to this package
func (p *Package) Uninstall() error {
	files, err := p.getSourceFiles()
	if err != nil {
		return err
	}

	files = append(files, p.metaDir, p.getDataDirectory())
	for _, fileName := range files {
		if err = os.RemoveAll(getAbsolutePath(fileName)); err != nil {
			return err
		}
	}

	return nil
}

// GetSize calculate all size of files in directories belonging to this package.
// Returns size in bytes
func (p *Package) GetSize() (int64, error) {
	files, err := p.getSourceFiles()
	if err != nil {
		return 0, err
	}

	files = append(files, p.metaDir, p.getDataDirectory())
	var size int64
	for _, fileName := range files {
		err = filepath.Walk(getAbsolutePath(fileName), func(_ string, info os.FileInfo, _ error) error {
			if info != nil && !info.IsDir() {
				size += info.Size()
			}
			return nil
		})
	}

	return size, nil
}

// Extra is used simultaneously for dependencies and extra packages
type Extra struct {
	Name       string
	Conditions []expression.Condition
	Compatible bool
}

// formatName formats the directory name to a single view
func formatName(dirName string) string {
	return strings.ToLower(strings.ReplaceAll(dirName, "-", "_"))
}

// parseFormat parse the directory name and returns its format.
// Example: "requests-2.26.0.dist-info" => "dist-info"
func parseFormat(dirName string) string {
	return filepath.Ext(dirName)[1:]
}

// clearVersion removes Extra characters from the version.
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
// pkgName. Returns the original directory name. Returns an error if the package
// is missing
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

	return "", ferror.PackageDirectoryMissing
}
