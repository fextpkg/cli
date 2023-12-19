package installer

import (
	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/pkg"
)

// Query is a struct for unifying the packages that need to be installed.
// It contains all the necessary parameters for searching in web repositions.
type Query struct {
	// Clean package name used for searching in the repository
	pkgName string
	// Conditions (operators and versions) to use for searching in the repository
	conditions []expression.Condition
	// Duplicate struct in case a search for extra packages was performed
	// without the required package already installed
	extraNames *Query
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

// copyQuery creates a new query similar to the one passed
func copyQuery(q *Query) *Query {
	return &Query{
		pkgName:    q.pkgName,
		conditions: q.conditions,
		extraNames: q.extraNames,
	}
}

// extraDependenciesToQuery converts the pkg.Dependency list to a Query list
func extraDependenciesToQuery(extras []pkg.Dependency) []*Query {
	var q []*Query
	for _, extra := range extras {
		q = append(q, newQuery(extra.PackageName, extra.Conditions))
	}

	return q
}
