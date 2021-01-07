package cfg

import (
	"github.com/go-ini/ini"
	"os"
)

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
	PathToConfigDir string

	ConfigFile *ini.File
)

func Load() {
	config, err := ini.Load(PathToConfigDir + CONFIG_FILE_NAME)
	if err != nil {
		// insert default cfg
		config = ini.Empty()
		_, err := config.Section("main").NewKey("libDir", "")
		if err != nil {
			panic(err)
		}
	}

	ConfigFile = config
}

func Save() {
	err := ConfigFile.SaveTo(PathToConfigDir + CONFIG_FILE_NAME)
	if err != nil {
		panic(err)
	}
}

