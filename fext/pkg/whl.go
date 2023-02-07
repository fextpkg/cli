package pkg

import (
	"github.com/fextpkg/cli/fext/utils"
)

// extra is used simultaneously for dependencies and extra packages
type extra struct {
	Name       string
	Conditions string
	markers    string
}

// CheckMarkers checks the possibility of installation according to the
// specified markers. Returns an error if parsing failed
func (e *extra) CheckMarkers() (bool, error) {
	// TODO move marker replaces from markers module to this func
	return utils.CompareExpression(e.markers)
}
