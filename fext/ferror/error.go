package ferror

import "errors"

var (
	SyntaxError             = errors.New("syntax error")
	MissingDirectory        = errors.New("package metadata directory not found")
	PackageAlreadyInstalled = errors.New("package already installed")
	NoSuitableVersion       = errors.New("no suitable version")
)

type MissingExtra struct {
	name string
}

func (e *MissingExtra) Error() string {
	return "extra not found: " + e.name
}

func NewMissingExtra(name string) error {
	return &MissingExtra{name: name}
}
