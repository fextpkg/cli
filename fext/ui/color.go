package ui

import (
	"fmt"
	"strings"
)

// It's lightweight module for colorize text
// NOTE! Function declares here will work correctly only in Go 1.10+,
// cause older versions not supported strings Builder.

const (
	PrefixCode = "\033["
	Reset      = PrefixCode + "0m"

	Red    = PrefixCode + "31m"
	Green  = PrefixCode + "32m"
	Orange = PrefixCode + "33m"

	Bold = PrefixCode + "1m"
)

func colorString(startColor, text string) string {
	b := strings.Builder{}
	b.WriteString(startColor)
	b.WriteString(text)
	b.WriteString(Reset)

	return b.String()
}

func PrintfOK(text string, args ...interface{}) (int, error) {
	return fmt.Print(colorString(Green+Bold, fmt.Sprintf(text, args...)))
}

func PrintfWarning(text string, args ...interface{}) (int, error) {
	return fmt.Print(colorString(Orange+Bold, fmt.Sprintf(text, args...)))
}

func PrintfError(text string, args ...interface{}) (int, error) {
	return fmt.Print(colorString(Red+Bold, fmt.Sprintf(text, args...)))
}

func PrintlnPlus(text string) (int, error) {
	b := strings.Builder{}
	b.WriteString(Green)
	b.WriteString(Bold)
	b.WriteString("+ ")
	b.WriteString(Reset)
	b.WriteString(text)

	return fmt.Println(b.String())
}

func PrintlnMinus(text string) (int, error) {
	b := strings.Builder{}
	b.WriteString(Red)
	b.WriteString(Bold)
	b.WriteString("- ")
	b.WriteString(Reset)
	b.WriteString(text)

	return fmt.Print(b.String())
}