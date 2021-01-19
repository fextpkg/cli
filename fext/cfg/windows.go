//+build windows

package cfg

import "os"

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	OS_LIB_PATH = homeDir + "\\AppData\\Local\\Programs\\Python\\"
}

var (
	OS_LIB_PATH string
)

const (
	PATH_TO_SITE_PACKAGES = "\\Lib\\site-packages"
	PYTHON_PREFIX = "Python"
	PYOS = "win32"
)
