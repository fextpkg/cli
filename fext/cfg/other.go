//+build !windows,!linux,!darwin

package cfg

import (
	"github.com/fextpkg/cli/fext/color"
	"runtime"
)

func init() {
	color.PrintfWarning("Fext doesn't support %s. This means that it doesn't guarantee stable operation",
						runtime.GOOS)
}

const (
	OS_LIB_PATH = "/"
	PATH_TO_SITE_PACKAGES = ""
	PYTHON_PREFIX = ""
	PYOS = runtime.GOOS
)
