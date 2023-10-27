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
	"log"
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	var originalMode uint32
	stdout := windows.Handle(os.Stdout.Fd())

	err := windows.GetConsoleMode(stdout, &originalMode)
	if err != nil {
		log.Fatal(err)
	}

	err = windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	if err != nil {
		log.Fatal(err)
	}
}
