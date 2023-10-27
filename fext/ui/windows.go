//go:build windows

/*
Please note: This file provides the ability to use the ANSI escape sequence only
for Windows 10 and above, as it has native support. By default, it is disabled.
The function in this file enables it until the terminal session is closed.

All versions below Windows 10 will not display colors in the console. To view
colorized strings, use another terminal like: Cmder, ConEmu, ANSICON or Mintty
*/

package ui

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	var originalMode uint32
	stdout := windows.Handle(os.Stdout.Fd())

	// Add modes only if it is not a redirect
	if err := windows.GetConsoleMode(stdout, &originalMode); err == nil {
		mode := originalMode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		// In case of an error, the changes are not applied. Continue the work of
		// application without the support of colorizing.
		_ = windows.SetConsoleMode(stdout, mode)
	}
}
