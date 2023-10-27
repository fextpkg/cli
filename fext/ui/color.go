package ui

import (
	"fmt"
	"os"
	"strings"
)

// This is a lightweight module that uses ANSI escape sequence to colorize text

const (
	PrefixCode = "\033["
	Reset      = PrefixCode + "0m"

	Red    = PrefixCode + "31m"
	Green  = PrefixCode + "32m"
	Orange = PrefixCode + "33m"

	Bold = PrefixCode + "1m"
)

// colorString colors the string in the specified color, using ANSI color escape
// sequences.
func colorString(startColor, text string) string {
	b := strings.Builder{}
	b.WriteString(startColor)
	b.WriteString(text)
	b.WriteString(Reset)

	return b.String()
}

// plusString generates a string with a green bold plus and your text.
func plusString(text string) string {
	b := strings.Builder{}
	b.WriteString(Green)
	b.WriteString(Bold)
	b.WriteString("+ ")
	b.WriteString(Reset)
	b.WriteString(text)

	return b.String()
}

// minusString generates a string with a red bold minus and your text.
func minusString(text string) string {
	b := strings.Builder{}
	b.WriteString(Red)
	b.WriteString(Bold)
	b.WriteString("- ")
	b.WriteString(Reset)
	b.WriteString(text)

	return b.String()
}

// PrintfOK generates a green bold string using fmt.Sprintf, and outputs it using
// fmt.Print.
func PrintfOK(text string, args ...interface{}) {
	fmt.Print(colorString(Green+Bold, fmt.Sprintf(text, args...)))
}

// PrintfWarning generates an orange bold string using fmt.Sprintf, and outputs
// it using fmt.Print.
func PrintfWarning(text string, args ...interface{}) {
	fmt.Print(colorString(Orange+Bold, fmt.Sprintf(text, args...)))
}

// PrintlnError generates a red bold string using the passed arguments and
// outputs it using the fmt.Println.
func PrintlnError(a ...string) {
	fmt.Println(colorString(Red+Bold, strings.Join(a, " ")))
}

// PrintfError generates a red bold string using fmt.Sprintf, and outputs it
// using fmt.Print.
func PrintfError(text string, args ...interface{}) {
	fmt.Print(colorString(Red+Bold, fmt.Sprintf(text, args...)))
}

// PrintlnPlus outputs a string with a green plus and specified text, using
// fmt.Println.
func PrintlnPlus(text string) {
	fmt.Println(plusString(text))
}

// PrintfPlus generates a string with a green plus using fmt.Sprint, and outputs
// it using fmt.Print.
func PrintfPlus(text string, args ...interface{}) {
	fmt.Print(plusString(fmt.Sprintf(text, args...)))
}

// PrintfMinus generates a string with a red minus using fmt.Sprint, and outputs
// it using fmt.Print.
func PrintfMinus(text string, args ...interface{}) {
	fmt.Print(minusString(fmt.Sprintf(text, args...)))
}

// Fatal is prints a text using PrintlnError and terminates the process with the
// status code 1. Works by analogy with log.Fatal, but uses colorized output.
func Fatal(a ...string) {
	PrintlnError(a...)
	os.Exit(1)
}
