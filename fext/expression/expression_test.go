package expression

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fextpkg/cli/fext/ferror"
)

var invalidExtraQueries = []string{
	"package[",
	"package]",
	"package[extraName",
	"package[extraName]]",
	"package[]",
	"package[[extraName]]",
}

func TestParseExtraNames(t *testing.T) {
	pkgName, extraName, err := ParseExtraNames("package[extraName]")
	assert.Nil(t, err)
	assert.Equal(t, extraName, []string{"extraName"})
	assert.Equal(t, pkgName, "package")

	pkgName, extraName, err = ParseExtraNames("package2[extraName2]>=1.0!=2.0")
	assert.Nil(t, err)
	assert.Equal(t, extraName, []string{"extraName2"})
	assert.Equal(t, pkgName, "package2>=1.0!=2.0")

	pkgName, extraName, err = ParseExtraNames("package3[extraName3, extraName4]>=1.0<2")
	assert.Nil(t, err)
	assert.Equal(t, extraName, []string{"extraName3", "extraName4"})
	assert.Equal(t, pkgName, "package3>=1.0<2")

	pkgName, extraName, err = ParseExtraNames("package4")
	assert.Nil(t, err)
	assert.Empty(t, extraName)
	assert.Equal(t, pkgName, "package4")
}

func TestParseExtraNamesInvalid(t *testing.T) {
	for _, query := range invalidExtraQueries {
		_, _, err := ParseExtraNames(query)
		assert.ErrorIs(t, err, ferror.SyntaxError)
	}
}
