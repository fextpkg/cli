package ui

import (
	"fmt"
	"os"
	"strings"
)

// It's lightweight module for colorize text

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

func plusString(text string) string {
	b := strings.Builder{}
	b.WriteString(Green)
	b.WriteString(Bold)
	b.WriteString("+ ")
	b.WriteString(Reset)
	b.WriteString(text)

	return b.String()
}

func minusString(text string) string {
	b := strings.Builder{}
	b.WriteString(Red)
	b.WriteString(Bold)
	b.WriteString("- ")
	b.WriteString(Reset)
	b.WriteString(text)

	return b.String()
}

func PrintfOK(text string, args ...interface{}) {
	fmt.Print(colorString(Green+Bold, fmt.Sprintf(text, args...)))
}

func PrintfWarning(text string, args ...interface{}) {
	fmt.Print(colorString(Orange+Bold, fmt.Sprintf(text, args...)))
}

func PrintlnError(a ...string) {
	fmt.Println(colorString(Red+Bold, strings.Join(a, " ")))
}

func PrintfError(text string, args ...interface{}) {
	fmt.Print(colorString(Red+Bold, fmt.Sprintf(text, args...)))
}

func PrintlnPlus(text string) {
	fmt.Println(plusString(text))
}

func PrintfPlus(text string, args ...interface{}) {
	fmt.Print(plusString(fmt.Sprintf(text, args...)))
}

func PrintfMinus(text string, args ...interface{}) {
	fmt.Print(minusString(fmt.Sprintf(text, args...)))
}

func Fatal(a ...string) {
	PrintlnError(a...)
	os.Exit(1)
}
