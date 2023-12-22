//go:build windows

package command

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ui"
)

type Debug struct{}

// DetectFlags does nothing and is a stub to maintain a single interface of
// interaction.
func (cmd *Debug) DetectFlags() error {
	return nil
}

// Execute prints debug info.
func (cmd *Debug) Execute() {
	fmt.Printf(
		"Fext (%s)\n\nLinked to: %s\nPython version: %s\nSystem platform: %s (tag: %s)\nChange mode: %v\nOS: %s, arch: %s\n",
		ui.BoldString(config.Version),
		ui.BoldString(config.PythonLibPath),
		ui.BoldString(config.PythonVersion),
		ui.BoldString(config.SysPlatform),
		ui.BoldString(config.PlatformTag),
		ui.BoldString(os.FileMode(config.DefaultChmod).String()),
		ui.BoldString(runtime.GOOS),
		ui.BoldString(runtime.GOARCH),
	)
}

// InitDebug initializes the "debug" command structure.
func InitDebug() *Debug {
	return &Debug{}
}
