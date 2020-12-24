package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// print beauty list of values and returns user choice
func ShowChoices(values []string) string {
	fmt.Println("Please choice one of next elements:")

	var elms = "\n"
	for i, value := range values {
		elms += strconv.Itoa(i+1) + ") " + value + "\t"
	}
	count := len(values)

	for {
		var i int
		fmt.Print(elms + "\n> ")
		_, err := fmt.Scanf("%d", &i)

		if err != nil {
			panic(err)
		}

		if i > 0 && i <= count {
			return values[i-1]
		} else {
			fmt.Println("< Element out of range!")
		}
	}
}

func ShowInput(text string) string {
	fmt.Println(text)
	var input, ok string

	for {
		fmt.Print("\n> ")
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			panic(err)
		}
		fmt.Print("Confirm written value? (y\\n) > ")
		_, err = fmt.Scanf("%s", &ok)
		if err != nil {
			panic(err)
		}

		if ok == "y" || ok == "Y" || ok == "yes" {
			break
		}
	}

	return input
}

// Parse first arguments and returns options and position where options end
func ParseOptions(args []string) ([]string, int) {
	var options []string
	for i, v := range args {
		if strings.HasPrefix(v, "-") {
			if strings.HasPrefix(v, "--") {
				v = v[2:]
			} else {
				v = v[1:]
			}

			options = append(options, v)
		} else {
			return options, i
		}
	}

	return options, 0
}

// clear last message, by set character at start position and send null strings
func ClearLastMessage(length int) {
	fmt.Printf("\r%s", strings.Repeat(" ", length))
}

func FindMinValue(values []int) int {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min
}

func ParseFormat(dir string) string {
	s := strings.Split(dir, ".")
	return s[len(s) - 1]
}

