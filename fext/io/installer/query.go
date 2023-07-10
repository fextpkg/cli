package installer

import (
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

// extrasToQuery converts the pkg.Extra list to a Query list
func extrasToQuery(extras []pkg.Extra) []*Query {
	var q []*Query
	for _, extraPackage := range extras {
		if extraPackage.Compatible {
			q = append(q, newQuery(extraPackage.Name, extraPackage.Conditions))
		}
	}

	return q
}
