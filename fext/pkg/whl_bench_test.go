package pkg

import (
	"testing"
)

var (
	pkgName = "sphinx"
	pkg, _  = Load(pkgName)
)

func BenchmarkPackage_parseMetaData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Load(pkgName)
	}
}

func BenchmarkPackage_GetExtraDependencies(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pkg.GetExtraDependencies("lint")
		if err != nil {
			return
		}
	}
}

func BenchmarkPackage_GetDependencies(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = pkg.GetDependencies()
	}
}

func BenchmarkPackage_GetSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := pkg.GetSize()
		if err != nil {
			return
		}
	}
}
