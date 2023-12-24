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
	// Name of directory that contains METADATA file
	metaDir string
	// Package name, as intended by the developer
	Name string
	// Package version
	Version string

	// A list of all dependencies specified through "Requires-Dist"
	Dependencies []Dependency
	// A list of all extra dependencies specified through "Provides-Extra"
	Extras []string
}

// Dependency used for both dependencies and extra dependencies simultaneously
type Dependency struct {
	// The raw, unprocessed string, stored exactly as it is in the metadata
	rawValue string
	// Processed Python markers, ready for comparison
	markers string
	// Determines if a dependency is an extra
	isExtra bool

	// Normalized package name obtained during the parsing of the raw string
	PackageName string
	// Normalized installation conditions for the required version, ready for comparison
	Conditions []expression.Condition
}

// Load a package by searching for its metadata directory and processing the
// metadata. Returns an error if it fails to process the metadata or if the
// directory is missing.
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

// LoadFromMetaDir parses the metadata in the specified path and returns a
// Package. Returns an error if any errors occur during the metadata processing.
func LoadFromMetaDir(path string) (*Package, error) {
	p := Package{
		metaDir:      path,
		Dependencies: []Dependency{},
		Extras:       []string{},
	}
	if err := p.parseMetaData(); err != nil {
		return nil, err
	}
	return &p, nil
}

// parseMetaData reads each line from the METADATA file in the specified
// metaDir. Sets the obtained values to public attributes.
// Returns an error if it fails to read the file.
func (p *Package) parseMetaData() error {
	data, err := os.ReadFile(getAbsolutePath(p.metaDir, "METADATA"))
	if err != nil {
		return err
	}

	// Metadata file is always separated by "\n"
	for _, s := range strings.Split(string(data), "\n") {
		// Remove unnecessary escape sequence characters
		s = strings.TrimSpace(s)
		// The file contains an empty line, which indicates that a long
		// description follows, but we don't need it
		if s != "" {
			field := strings.SplitN(s, ": ", 2)
			if len(field) != 2 {
				// Handle extreme situations when file has been poorly
				// generated
				continue
			}
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
// packages and modules. Returns its content. If the file doesn't exist, it
// returns a slice consisting of a single element - the name of the package
// itself.
func (p *Package) getTopLevel() ([]string, error) {
	files, err := io.ReadLines(getAbsolutePath(p.metaDir, "top_level.txt"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		// Add the package name manually, since some generators do not create a
		// "top_level.txt" file
		files = []string{formatName(p.Name)}
	}
	return files, nil
}

// getSourceFiles checks for the existence of directories and files obtained
// from getTopLevel. Returns only the existing files and directories belongs
// to this package.
func (p *Package) getSourceFiles() ([]string, error) {
	files, err := p.getTopLevel()
	if err != nil {
		return nil, err
	}
	for i, fileName := range files {
		// The names in the top_level.txt are stored without extensions.
		// This means that first we need to check if the name is a directory.
		// If not, it is a Python file
		if _, err = os.Stat(getAbsolutePath(fileName)); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, err
			}
			// Add the suffix ".py" and replace the current name
			files[i] = fileName + ".py"
		}
	}
	return files, nil
}

// getDataDirectory returns the name of the directory with the data files.
// It doesn't check if the directory exists.
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
			// Important: RemoveAll doesn't return an error if the file
			// not exists
			return err
		}
	}

	return nil
}

// GetSize calculates the total size of all files belonging to the package and
// returns the size in bytes.
func (p *Package) GetSize() (int64, error) {
	files, err := p.getSourceFiles()
	if err != nil {
		return 0, err
	}

	files = append(files, p.metaDir, p.getDataDirectory())
	var size int64
	for _, fileName := range files {
		_ = filepath.Walk(
			getAbsolutePath(fileName), func(_ string, info os.FileInfo, _ error) error {
				// Ignore the error as the data directory may not exist
				if info != nil && !info.IsDir() {
					// The weight of folders is always incorrect,
					// so calculate the weight based on all the files instead
					size += info.Size()
				}
				return nil
			},
		)
	}

	return size, nil
}

// GetMetaDirectoryPath returns the absolute path to the package's meta-directory.
func (p *Package) GetMetaDirectoryPath() string {
	return getAbsolutePath(p.metaDir)
}

// GetDependencies retrieves compatible package dependencies that are ready for
// comparison.
// Returns an error if there are any issues during metadata parsing.
func (p *Package) GetDependencies() ([]Dependency, error) {
	var packages []Dependency

	for _, dep := range p.Dependencies {
		if !dep.isExtra {
			// Markers may not always be present in dependency line
			if dep.markers != "" {
				compatible, err := expression.CompareMarkers(dep.markers)
				if err != nil {
					return nil, err
				} else if !compatible {
					continue
				}
			}

			dep.PackageName, dep.Conditions = expression.ParseConditions(dep.rawValue)
			packages = append(packages, dep)
		}
	}

	return packages, nil
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

// HasExtraName checks if the extra dependency name exists.
func (p *Package) HasExtraName(name string) bool {
	for _, depName := range p.Extras {
		if depName == name {
			return true
		}
	}

	return false
}

// formatName standardizes a directory name for easier searching among other
// directories and files. It replaces all "-" to "_" and converts the string to
// lowercase.
func formatName(dirName string) string {
	return strings.ToLower(strings.ReplaceAll(dirName, "-", "_"))
}

// parseExtension separates the directory/file extension.
// Example: "requests-2.26.0.dist-info" => "dist-info"
func parseExtension(dirName string) string {
	return filepath.Ext(dirName)[1:]
}

// clearVersion removes extra characters from the version.
// Example: "2.26.0.dist" => "2.26.0"
func clearVersion(version string) string {
	return strings.Replace(version, ".dist", "", 1)
}

// getAbsolutePath returns absolute path to the file in directory with packages
func getAbsolutePath(elem ...string) string {
	return config.PythonLibPath + string(os.PathSeparator) + filepath.Join(elem...)
}

// Parse directory by format "%pkgName%-%version%.%extension%" and returns it
func parseDirectoryName(dirName string) (string, string, string) {
	// [name, version, extension]
	meta := strings.SplitN(dirName, "-", 3)

	// avoid errors
	if len(meta) >= 2 {
		return meta[0], clearVersion(meta[1]), parseExtension(dirName)
	} else {
		return meta[0], "", ""
	}
}

// getPackageMetaDir searches for a folder with wheel extension from the
// specified pkgName. Returns the original directory name or an error if the
// package is missing
func getPackageMetaDir(pkgName string) (string, error) {
	dirInfo, err := os.ReadDir(config.PythonLibPath)
	if err != nil {
		return "", err
	}
	pkgName = formatName(pkgName)

	for _, dir := range dirInfo {
		curPkgName, _, ext := parseDirectoryName(dir.Name())
		if formatName(curPkgName) == pkgName && ext == "dist-info" {
			return dir.Name(), nil
		}
	}

	return "", ferror.PackageDirectoryMissing
}

// parseRequirement separates the name from the markers and determines whether
// the dependency if considered as extra.
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
