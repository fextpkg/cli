package cfg

import "os"

const (
	VERSION = "0.0.3"
	CONFIG_FILE_NAME = "fext"
	BASE_PACKAGE_URL = "https://pypi.org/simple/"
	MAX_MESSAGE_LENGTH = 30 // max letters in progress bar
	DEFAULT_CHMOD = 0775
	PATH_SEPARATOR = string(os.PathSeparator)
)

var ( // will fill in future, when program will runs
	PathToLib string
	PythonVersion string
)

