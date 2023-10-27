package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/ferror"
	"github.com/fextpkg/cli/fext/pkg"
	"github.com/fextpkg/cli/fext/ui"
)

type Freeze struct {
	// Print style mode
	printMode string

	// Directory names with packages metadata for scanning and output
	metaDirectories []string
}

// getPrintFunc selects the functions responsible for the specified print printMode.
// If there is not such a function, an error will be returned.
func (cmd *Freeze) getPrintFunc(mode string) (func(), error) {
	if mode == "human" {
		return cmd.printStyleHuman, nil
	} else {
		return nil, &ferror.UnexpectedMode{Mode: mode}
	}
}

// printStyleHuman prints a list of packages in a human-readable style, including
// general info about the number of packages and their weight.
//
// Function ignores packages with an error, which means the error will never be
// returned. Packages that failed to load will be displayed at the end as a
// warning line.
func (cmd *Freeze) printStyleHuman() {
	var brokenPackages []string
	var count int
	var size int64

	for _, metaDir := range cmd.metaDirectories {
		p, err := pkg.LoadFromMetaDir(metaDir)
		if err == nil {
			s, err := p.GetSize()
			if err != nil {
				brokenPackages = append(brokenPackages, metaDir)
				continue
			}

			size += s
			count++
			fmt.Printf("%s (%s)\n", p.Name, p.Version)
		} else {
			brokenPackages = append(brokenPackages, metaDir)
		}
	}

	fmt.Printf("\nTotal: %d (%.2f MB)\n", count, float32(size/1024)/1024)
	if len(brokenPackages) > 0 {
		ui.PrintfWarning(
			"During the package analysis, some errors has occurred with packages: %s\n",
			strings.Join(brokenPackages, ", "),
		)
	}
}

// scanMetaDirectories goes through the directory with python modules and
// packages, selects the meta-directories and appends them to the metaDirectories
// attribute. Returns an error if the folder could not be read.
func (cmd *Freeze) scanMetaDirectories() error {
	files, err := os.ReadDir(config.PythonLibPath)
	if err != nil {
		return err
	}

	// Go through the files and select meta directories (wheel has the "dist-info"
	// suffix)
	for _, f := range files {
		dirName := f.Name()
		if f.IsDir() && strings.HasSuffix(dirName, "dist-info") {
			cmd.metaDirectories = append(cmd.metaDirectories, dirName)
		}
	}

	return nil
}

// DetectFlags analyzes the passed flags and fills in the variables associated
// with them.
//
// Returns ferror.HelpFlag if you need to print the docstring about this command.
// Returns ferror.UnknownFlag if passed the unknown flag.
// Returns ferror.MissingOptionValue if the correct but empty option is passed.
func (cmd *Freeze) DetectFlags() error {
	for _, f := range config.Flags {
		if f == "h" || f == "help" { // Help flag
			return ferror.HelpFlag
		} else if strings.HasPrefix(f, "m=") || strings.HasPrefix(f, "printMode=") {
			// Print style printMode. To preserve performance and foundation, we do not add a
			// structure for flags with a value. Instead, parse these flag in only "freeze"
			// command. But it may have to be redone in the future. Who knows…
			cmd.printMode = strings.SplitN(f, "=", 2)[1]
			if cmd.printMode == "" { // Empty value
				return &ferror.MissingOptionValue{Opt: f[:len(f)-1]}
			}
		} else if f == "m" || f == "printMode" { // Passed "printMode" option without value
			return &ferror.MissingOptionValue{Opt: f}
		} else { // Unexpected flag
			return &ferror.UnknownFlag{Flag: f}
		}
	}

	return nil
}

// Execute prints a list of packages in the selected printMode.
func (cmd *Freeze) Execute() {
	fn, err := cmd.getPrintFunc(cmd.printMode)
	if err != nil {
		ui.Fatal("Unable to select print mode:", err.Error())
	}

	if err = cmd.scanMetaDirectories(); err != nil {
		ui.Fatal("Unable to scan meta directories:", err.Error())
	}

	fn()
}

// InitFreeze initializes the "freeze" command structure with default parameters.
func InitFreeze() *Freeze {
	return &Freeze{
		printMode: "human",
	}
}
