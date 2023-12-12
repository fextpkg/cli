package ferror

import "errors"

var (
	SyntaxError             = errors.New("syntax error")
	PackageDirectoryMissing = errors.New("package metadata directory not found")
	PackageAlreadyInstalled = errors.New("package already installed")
	NoSuitableVersion       = errors.New("no suitable version")
	HelpFlag                = errors.New("help flag")
	UnexpectedCommand       = errors.New("unexpected command")
)

type MissingExtra struct {
	Name string
}

func (e *MissingExtra) Error() string {
	return "extra not found: " + e.Name
}

type UnexpectedMode struct {
	Mode string
}

func (e *UnexpectedMode) Error() string {
	return "unexpected mode: " + e.Mode
}

type UnknownFlag struct {
	Flag string
}

func (e *UnknownFlag) Error() string {
	return "unknown flag: " + e.Flag
}

type MissingOptionValue struct {
	Opt string
}

func (e *MissingOptionValue) Error() string {
	return "option '" + e.Opt + "': missing value"
}

type UnexpectedMarker struct {
	Marker string
}

func (e *UnexpectedMarker) Error() string {
	return "unexpected marker: " + e.Marker
}

type UnexpectedOperator struct {
	Operator string
}

func (u *UnexpectedOperator) Error() string {
	return "unexpected operator: " + u.Operator
}
