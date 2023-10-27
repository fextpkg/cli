package command

import (
	"fmt"
	"runtime"

	"github.com/fextpkg/cli/fext/config"
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
		"Fext (%s)\n\nLinked to: %s\nOS: %s, arch: %s\n",
		config.Version,
		config.PythonLibPath,
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// InitDebug initializes the "debug" command structure.
func InitDebug() *Debug {
	return &Debug{}
}
