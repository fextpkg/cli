package main

import (
	"github.com/Flacy/fext/fext/cfg"
	"github.com/Flacy/fext/fext/cmd"
	"github.com/Flacy/fext/fext/help"
	"github.com/Flacy/fext/fext/utils"

	"fmt"
	"os"

	"github.com/go-ini/ini"
)

func loadConfig(configDir string) *ini.File {
	config, err := ini.Load(configDir + cfg.CONFIG_FILE_NAME)
	if err != nil {
		// insert default cfg
		config = ini.Empty()
		config.Section("main").NewKey("libDir", "")
	}

	return config
}

func saveConfig(configDir string, config *ini.File) {
	err := config.SaveTo(configDir + cfg.CONFIG_FILE_NAME)

	if err != nil {
		panic(err)
	}
}

func main() {
	// first argument it's name of executable file
	args := os.Args[1:]

	if len(args) == 0 {
		help.Show()
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		} else {
			configDir += cfg.PATH_SEPARATOR
		}

		command := args[0]
		args = args[1:]
		config := loadConfig(configDir)
		libDirKey := config.Section("main").Key("libDir")

		if libDirKey.String() == "" {
			// init
			libDirKey.SetValue(utils.GetPythonLibDirectory())
		}
		// do it, cause on windows separator doesn't saves
		cfg.PathToLib = libDirKey.Value() + cfg.PATH_SEPARATOR
		cfg.PythonVersion = utils.ParsePythonVersion(cfg.PathToLib)

		switch command {
		case "install", "i":
			cmd.Install(args)
		case "uninstall", "u":
			cmd.Uninstall(args)
		case "freeze":
			cmd.Freeze()
		case "debug":
			help.ShowDebug(configDir + cfg.CONFIG_FILE_NAME)
		default:
			fmt.Println("Unexpected command")
		}

		saveConfig(configDir, config)
	}
}
