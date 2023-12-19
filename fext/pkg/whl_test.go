package pkg

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/ui"
)

func formatMetaDirectory(name string) string {
	return fmt.Sprintf("%s-%s.dist-info", name, PackageVersion)
}

func createMetaDirectory(name string) (string, error) {
	metadataDir := formatMetaDirectory(name)
	pathToMetaDir := getAbsolutePath(metadataDir)
	return pathToMetaDir, os.Mkdir(pathToMetaDir, config.DefaultChmod)
}

func createTestPackage() error {
	pathToMetaDir, err := createMetaDirectory(PackageName)
	if err != nil {
		return err
	}

	pathToMetaFile := path.Join(pathToMetaDir, "METADATA")
	err = os.WriteFile(pathToMetaFile, []byte(Metadata), config.DefaultChmod)
	if err != nil {
		return err
	}

	return nil
}

func createBrokenPackage() (string, error) {
	name := "test_broken"
	_, err := createMetaDirectory(name)
	return name, err
}

func cleanUpPackage(name string) {
	err := os.RemoveAll(getAbsolutePath(formatMetaDirectory(name)))
	if err != nil {
		ui.Fatal("Unable to cleanup package: " + name + ": " + err.Error())
	}
}

func init() {
	p, err := Load(PackageName)
	if err == nil {
		if err = os.RemoveAll(p.GetMetaDirectoryPath()); err != nil {
			ui.Fatalf(
				"Unable to initiate test package. There is a conflict with an existing package: %s: %s",
				PackageName,
				err,
			)
		}
	}

	err = createTestPackage()
	if err != nil {
		ui.Fatal("Unable to initiate test package: " + err.Error())
	}
}

func containsItem(items []string, item string) bool {
	for _, dep := range items {
		if dep == item {
			return true
		}
	}

	return false
}

func TestPackage_GetExtraPackages(t *testing.T) {
	p, err := Load(PackageName)
	assert.Nil(t, err)

	for _, extraName := range PackageExtras {
		deps, err := p.GetExtraDependencies(extraName)
		assert.Nil(t, err)
		assert.Len(t, deps, 1)

		for _, dep := range deps {
			assert.True(t, dep.isExtra)
			assert.NotEmpty(t, dep.markers)
		}
	}

	for _, extraName := range PackageExtrasUnknown {
		deps, err := p.GetExtraDependencies(extraName)
		assert.Nil(t, err)
		assert.Empty(t, deps)
	}
}

func TestPackage_GetDependencies(t *testing.T) {
	p, err := Load(PackageName)
	assert.Nil(t, err)

	deps := p.GetDependencies()
	for _, dep := range deps {
		assert.False(t, dep.isExtra)
		assert.Empty(t, dep.markers)
		assert.True(t, containsItem(PackageDependencies, dep.PackageName))
	}
}

func TestPackage_GetSize(t *testing.T) {
	p, err := Load(PackageName)
	assert.Nil(t, err)

	size, err := p.GetSize()
	assert.Nil(t, err)
	assert.NotZero(t, size)
}

func TestPackage_getTopLevel(t *testing.T) {
	p, err := Load(PackageName)
	assert.Nil(t, err)

	topLevel, err := p.getTopLevel()
	assert.Nil(t, err)
	assert.Len(t, topLevel, 1)
	assert.Equal(t, topLevel, []string{PackageName})
}

func TestPackage_getSourceFiles(t *testing.T) {
	p, err := Load(PackageName)
	assert.Nil(t, err)

	sourceFiles, err := p.getSourceFiles()
	assert.Nil(t, err)
	// Since we did not create a directory with source files, the algorithm
	// assumes that this package is a module (a .py file)
	assert.Equal(t, sourceFiles, []string{PackageName + ".py"})
}

func TestPackage_Load(t *testing.T) {
	p, err := Load(PackageName)
	assert.Nil(t, err)

	match, err := expression.CompareVersion(p.Version, "==", PackageVersion)
	assert.Nil(t, err)
	assert.True(t, match)

	for _, provideExtra := range p.Extras {
		if provideExtra == "security" {
			// It remains unclear what to do with such packages until the end.
			// In normal code, the algorithm would throw an error stating that
			// there are no extra packages because it only checks the
			// 'Requires-Dist' key. However, the extra packages might be
			// specified in the metadata using the 'Provide-Extra' key.
			continue
		}

		assert.True(t, containsItem(PackageExtras, provideExtra))
	}

	metaDir := fmt.Sprintf("%s-%s.dist-info", PackageName, PackageVersion)

	assert.Equal(t, p.metaDir, metaDir)
	assert.Equal(t, p.GetMetaDirectoryPath(), filepath.Join(config.PythonLibPath, metaDir))

	// cleanup test package
	err = p.Uninstall()
	assert.Nil(t, err)
}

func TestPackage_LoadFromMetaDir(t *testing.T) {
	p, err := Load(PackageName)
	assert.Nil(t, err)

	pFromMeta, err := LoadFromMetaDir(p.metaDir)
	assert.Nil(t, err)

	assert.Equal(t, p.Name, pFromMeta.Name)
	assert.Equal(t, p.Version, pFromMeta.Version)
	assert.Equal(t, p.GetDependencies(), pFromMeta.GetDependencies())

	s, err := p.GetSize()
	s2, err := pFromMeta.GetSize()
	assert.Equal(t, s, s2)
}

func TestPackage_LoadMissing(t *testing.T) {
	_, err := Load("something")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ferror.PackageDirectoryMissing)
}

func TestPackage_ParseMetaDataError(t *testing.T) {
	name, err := createBrokenPackage()
	assert.Nil(t, err)

	_, err = Load(name)
	assert.NotNil(t, err)

	_, err = LoadFromMetaDir(formatMetaDirectory(name))
	assert.NotNil(t, err)

	t.Cleanup(func() { cleanUpPackage(name) })
}

func TestGetPackageMetaDirFail(t *testing.T) {
	_, err := getPackageMetaDir("some_not_exists")
	assert.NotNil(t, err)

	cleanUpPackage(PackageName)
	err = createTestPackage()
	assert.Nil(t, err)

	normalPackage, err := Load(PackageName)
	assert.Nil(t, err)

	// change for tests
	pythonLibPath := config.PythonLibPath
	config.PythonLibPath = "//"

	// tests
	_, err = getPackageMetaDir(PackageName)
	assert.NotNil(t, err)
	_, err = Load(PackageName)
	assert.NotNil(t, err)
	_, err = normalPackage.getSourceFiles()
	assert.NotNil(t, err)
	_, err = normalPackage.getTopLevel()
	assert.NotNil(t, err)
	err = normalPackage.parseMetaData()
	assert.NotNil(t, err)
	size, err := normalPackage.GetSize()
	assert.NotNil(t, err)
	assert.Zero(t, size)
	err = normalPackage.Uninstall()
	assert.NotNil(t, err)

	// restore
	config.PythonLibPath = pythonLibPath

	t.Cleanup(func() { cleanUpPackage(PackageName) })
}

func TestBrokenMetadata(t *testing.T) {
	name, err := createBrokenPackage()
	assert.Nil(t, err)

	metaDirectoryPath := getAbsolutePath(formatMetaDirectory(name))
	metadataFilePath := filepath.Join(metaDirectoryPath, "METADATA")
	err = os.WriteFile(metadataFilePath, []byte(MetadataBrokenMarker), config.DefaultChmod)
	assert.Nil(t, err)

	p, err := Load(name)
	assert.Nil(t, err)

	var expectedError *ferror.UnexpectedMarker
	deps, err := p.GetExtraDependencies("socks")
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &expectedError)
	assert.Empty(t, deps)

	err = os.WriteFile(metadataFilePath, []byte(MetadataInvalidSyntaxMarker), config.DefaultChmod)
	assert.Nil(t, err)

	p, err = Load(name)
	assert.Nil(t, err)

	deps, err = p.GetExtraDependencies("socks")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ferror.SyntaxError)
	assert.Zero(t, deps)

	t.Cleanup(func() { cleanUpPackage(name) })
}
