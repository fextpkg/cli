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

	Extras []string

	Dependencies []Dependency
}

// Dependency used for both dependencies and extra packages simultaneously
type Dependency struct {
	rawValue string
	markers  string
	isExtra  bool

	// Attributes to be filled during processing and comparison of rawValue
	PackageName string
	Conditions  []expression.Condition
}

func Load(pkgName string) (*Package, error) {
	dirName, err := getPackageMetaDir(pkgName)
	if err != nil {
		return nil, err
	}

	p := Package{
		Name:         pkgName,
		metaDir:      dirName,
		Dependencies: []Dependency{},
		Extras:       []string{},
	}
	if err = p.parseMetaData(); err != nil {
		return nil, err
	}

	return &p, nil
}

func LoadFromMetaDir(metaDir string) (*Package, error) {
	p := Package{
		metaDir:      metaDir,
		Dependencies: []Dependency{},
		Extras:       []string{},
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

	for _, s := range strings.Split(string(data), "\n") {
		s = strings.TrimSpace(s)
		if s != "" {
			field := strings.SplitN(s, ": ", 2)
			key, value := field[0], field[1]

			switch key {
			case "Requires-Dist":
				p.Dependencies = append(p.Dependencies, parseRequirement(value))
			case "Provides-Extra":
				p.Extras = append(p.Extras, value)
			case "Version":
				p.Version = value
			case "Name":
				p.Name = value
			}
		} else {
			break
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
		files = []string{formatName(p.Name)}
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

// GetMetaDirectoryPath gets absolute path to the package meta directory.
func (p *Package) GetMetaDirectoryPath() string {
	return getAbsolutePath(p.metaDir)
}

func (p *Package) GetDependencies() []Dependency {
	var packages []Dependency

	for _, dep := range p.Dependencies {
		if !dep.isExtra {
			dep.PackageName, dep.Conditions = expression.ParseConditions(dep.rawValue)
			packages = append(packages, dep)
		}
	}

	return packages
}

// GetExtraDependencies retrieves compatible extra dependencies and returns an
// empty slice if none are found.
// Returns an error if there are any issues during metadata parsing.
func (p *Package) GetExtraDependencies(extraName string) ([]Dependency, error) {
	var extraPackages []Dependency

	for _, dep := range p.Dependencies {
		if dep.isExtra {
			// Initially, verify if the expression has the "extra" marker
			// with the required value
			match, err := expression.MatchExtraMarker(dep.markers, extraName)
			if err != nil {
				return nil, err
			} else if match {
				// After a successful search, validate the entire expression for
				// truth, considering that some extra dependencies may include
				// additional markers themselves
				compatible, err := expression.CompareMarkers(dep.markers)
				if err != nil {
					return nil, err
				} else if compatible {
					// Lastly, parse the conditions and extract the
					// package name for further handling
					dep.PackageName, dep.Conditions = expression.ParseConditions(dep.rawValue)
					extraPackages = append(extraPackages, dep)
				}
			}
		}
	}

	return extraPackages, nil
}

func (p *Package) HasExtraName(name string) bool {
	for _, depName := range p.Extras {
		if depName == name {
			return true
		}
	}

	return false
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
	return config.PythonLibPath + string(os.PathSeparator) + filepath.Join(elem...)
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

func parseRequirement(s string) Dependency {
	var markers string

	// [name_and_conditions, markers]
	exp := strings.SplitN(s, "; ", 2)
	if len(exp) > 1 { // Has markers
		markers = exp[1]
	} else { // No markers provided
		markers = ""
	}

	return Dependency{
		rawValue: exp[0],
		markers:  markers,
		isExtra:  strings.Contains(markers, "extra"),
	}
}
