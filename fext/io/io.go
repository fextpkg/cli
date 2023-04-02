package io

import (
	"bufio"
	"os"
	"strings"
)

func ReadLines(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(f)
	var result []string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line != "" {
			result = append(result, line)
		}
	}

	return result, f.Close()
}
