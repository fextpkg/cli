package ferror

import "errors"

var (
	// SyntaxError means that the passed query can't be processed due to syntax errors
	SyntaxError = errors.New("syntax error")
	// PackageDirectoryMissing means that the directory containing the package
	// metadata was not found. The directory with the .dist-info extension is missing.
	PackageDirectoryMissing = errors.New("package metadata directory not found")
	// PackageAlreadyInstalled means that you are trying to install a package
	// already installed in the system.
	PackageAlreadyInstalled = errors.New("package already installed")
	// PackageInLocalList means that the package installation attempt is being
	// made for the second time. The previous suitable version has already been
	// installed. The error arises during the comparison of the installed
	// package version with the operators obtained from other dependent packages.
	// To summarize, if this error occurs after comparing the list of operators,
	// it means that reinstalling the package is not necessary. It is already
	// present in the system and compatible with the other packages that use it.
	PackageInLocalList = errors.New("package version compared and already installed")
	// NoSuitableVersion means that no suitable version was found for the given query.
	// This can be related to both the operators provided and the package's
	// complete incompatibility with the current system (e.g., different Python
	// version or unsupported platform).
	NoSuitableVersion = errors.New("no suitable version")
	// HelpFlag means that the help string for the given command needs to be
	// displayed on the screen.
	HelpFlag = errors.New("help flag")
	// UnexpectedCommand means that it was not possible to determine the
	// command that was passed.
	UnexpectedCommand = errors.New("unexpected command")
)

// MissingExtra means that the extra package names provided were not found.
type MissingExtra struct {
	Name string
}

func (e *MissingExtra) Error() string {
	return "extra not found: " + e.Name
}

// UnexpectedMode means that an unknown print mode was passed for the "freeze" command.
type UnexpectedMode struct {
	Mode string
}

func (e *UnexpectedMode) Error() string {
	return "unexpected mode: " + e.Mode
}

// UnknownFlag means that an unexpected flag was passed for the given command.
type UnknownFlag struct {
	Flag string
}

func (e *UnknownFlag) Error() string {
	return "unknown flag: " + e.Flag
}

// MissingOptionValue means that the flag requires a value, but it was found
// to be empty.
type MissingOptionValue struct {
	Opt string
}

func (e *MissingOptionValue) Error() string {
	return "option '" + e.Opt + "': missing value"
}

// UnexpectedMarker means that an unknown marker was passed.
// All markers must follow the PEP 508 standard.
type UnexpectedMarker struct {
	Marker string
}

func (e *UnexpectedMarker) Error() string {
	return "unexpected marker: " + e.Marker
}

// UnexpectedOperator means that an unknown comparison/logical operator was passed.
type UnexpectedOperator struct {
	Operator string
}

func (u *UnexpectedOperator) Error() string {
	return "unexpected operator: " + u.Operator
}
