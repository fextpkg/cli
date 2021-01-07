package main

import (
	"github.com/Flacy/fext/fext/cfg"
	"github.com/Flacy/fext/fext/cmd"
	"github.com/Flacy/fext/fext/help"
	"github.com/Flacy/fext/fext/utils"

	"fmt"
	"os"
)

func initBaseVariables() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfg.PathToConfigDir = configDir + cfg.PATH_SEPARATOR
}

func initVariables() {
	libDirKey := cfg.ConfigFile.Section("main").Key("libDir")
	if libDirKey.String() == "" {
		// init
		libDirKey.SetValue(utils.GetPythonLibDirectory())
	}
	// do it, cause on windows separator doesn't saves
	cfg.PathToLib = libDirKey.Value() + cfg.PATH_SEPARATOR
	cfg.PythonVersion = utils.ParsePythonVersion(cfg.PathToLib)
}

func initConfig() {
	initBaseVariables()
	cfg.Load()
	initVariables()
}

func main() {
	// first argument it's name of executable file
	args := os.Args[1:]

	if len(args) == 0 {
		help.Show()
	} else {
		command := args[0]
		args = args[1:]

		initConfig()

		switch command {
		case "install", "i":
			cmd.Install(args)
		case "uninstall", "u":
			cmd.Uninstall(args)
		case "freeze":
			cmd.Freeze()
		case "debug":
			help.ShowDebug(cfg.PathToConfigDir + cfg.CONFIG_FILE_NAME)
		default:
			fmt.Println("Unexpected command")
		}

		cfg.Save()
	}
}

