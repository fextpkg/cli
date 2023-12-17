package installer

import (
	"fmt"

	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/pkg"
)

// Query is auxiliary struct for unification of packages to be installed
type Query struct {
	pkgName    string
	conditions []expression.Condition
}

// newRawQuery pre-parses conditional statements and creates a new query
func newRawQuery(s string) *Query {
	pkgName, conditions := expression.ParseConditions(s)
	fmt.Println(">>", s, pkgName, conditions)
	return &Query{
		pkgName:    pkgName,
		conditions: conditions,
	}
}

// newQuery creates a new query with already known parameters
func newQuery(pkgName string, conditions []expression.Condition) *Query {
	return &Query{
		pkgName:    pkgName,
		conditions: conditions,
	}
}

// extrasToQuery converts the pkg.Dependency list to a Query list
func extrasToQuery(extras []pkg.Dependency) []*Query {
	var q []*Query
	for _, extraPackage := range extras {
		q = append(q, newQuery(extraPackage.PackageName, extraPackage.Conditions))
	}

	return q
}
