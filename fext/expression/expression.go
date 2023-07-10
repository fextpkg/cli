package expression

import (
	"strings"

	"github.com/fextpkg/cli/fext/ferror"
)

// ParseExtraNames separates extra dependencies on the name and conditions of the
// package (PEP 685). Returns pkgName, extraNames and ferror.SyntaxError if
// syntax is invalid.
//
//	ParseExtraNames("package[extra]>=1")
//	=> "package>=1", ["extra"], nil
func ParseExtraNames(s string) (string, []string, error) {
	startQuote := strings.Index(s, "[")
	endQuote := strings.LastIndex(s, "]")

	if startQuote != -1 && endQuote != -1 {
		originalName := s[:startQuote] + s[endQuote+1:] // pkgName and conditions
		s = s[startQuote+1 : endQuote]                  // extra names
		if strings.ContainsAny(s, "[]") {
			return originalName, nil, ferror.SyntaxError
		}

		var extraNames []string
		s = strings.ReplaceAll(s, " ", "")
		for _, name := range strings.Split(s, ",") {
			extraNames = append(extraNames, name)
		}

		return originalName, extraNames, nil
	} else if startQuote != -1 || endQuote != -1 {
		return s, nil, ferror.SyntaxError
	}

	return s, nil, nil
}
