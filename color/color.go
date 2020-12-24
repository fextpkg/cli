package color

import (
	"fmt"
	"strings"
)

// It's lightweight module for colorize text
// NOTE! Function declares here will work correctly only in Go 1.10+,
// cause older versions not supported strings Builder.

const (
	PREFIX_CODE = "\033["
	RESET = PREFIX_CODE + "0m"

	// colors
	RED = PREFIX_CODE + "31m"
	GREEN = PREFIX_CODE + "32m"

	// effects
	BOLD = PREFIX_CODE + "1m"
)

func initBuilder() strings.Builder {
	return strings.Builder{}
}

func getColoredString(startColor, text string) string {
	b := initBuilder()
	b.WriteString(startColor)
	b.WriteString(text)
	b.WriteString(RESET)

	return b.String()
}

func PrintfError(text string, args ...interface{}) (int, error) {
	coloredString := getColoredString(RED + BOLD, fmt.Sprintf(text, args...))
	return fmt.Print(coloredString)
}

func PrintfOK(text string, args ...interface{}) (int, error) {
	coloredString := getColoredString(GREEN + BOLD, fmt.Sprintf(text, args...))
	return fmt.Print(coloredString)
}

func PrintflnStatus(text, color, status string, args ...interface{}) (int, error) {
	b := initBuilder()
	b.WriteString(fmt.Sprintf(text + " ", args...))
	b.WriteString(color)
	b.WriteString("(" + status + ")")
	b.WriteString(RESET)

	return fmt.Println(b.String())
}

func PrintflnStatusOK(text string, args ...interface{}) (int, error) {
	return PrintflnStatus(text, GREEN + BOLD, "OK", args...)
}

func PrintflnStatusError(text, reason string, args ...interface{}) (int, error) {
	return PrintflnStatus(text, RED + BOLD, reason, args...)
}
