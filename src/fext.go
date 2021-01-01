package main

import (
	"github.com/Flacy/fext/base_cfg"
	"github.com/Flacy/fext/cmd"
	"github.com/Flacy/fext/help"
	"github.com/Flacy/fext/utils"

	"fmt"
	"os"

	"github.com/go-ini/ini"
)

func loadConfig(configDir string) *ini.File {
	config, err := ini.Load(configDir + base_cfg.CONFIG_NAME)
	if err != nil {
		// insert default cfg
		config = ini.Empty()
		config.Section("main").NewKey("libDir", "")
	}

	return config
}

func saveConfig(configDir string, config *ini.File) {
	err := config.SaveTo(configDir + base_cfg.CONFIG_NAME)

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
			configDir += string(os.PathSeparator)
		}

		command := args[0]
		args = args[1:]
		config := loadConfig(configDir)
		libDirKey := config.Section("main").Key("libDir")

		if libDirKey.String() == "" {
			// init
			libDirKey.SetValue(utils.FindPythonLibDirectory())
		}
		// do it, cause on windows separator doesn't saves
		libDir := libDirKey.Value() + string(os.PathSeparator)

		switch command {
		case "install", "i":
			cmd.Install(libDir, args)
		case "uninstall", "u":
			cmd.Uninstall(libDir, args)
		case "freeze":
			cmd.Freeze(libDir)
		case "debug":
			help.ShowDebug(configDir +base_cfg.CONFIG_NAME, libDir)
		default:
			fmt.Println("Unexpected command")
		}

		saveConfig(configDir, config)
	}
}
