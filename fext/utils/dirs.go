package utils

import (
	"github.com/fextpkg/cli/fext/cfg"
	"github.com/fextpkg/cli/fext/color"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const OS = runtime.GOOS

func SelectOneDirectory(baseDir string, versions []string) string {
	var path string
	countV := len(versions)
	if countV == 0 {
		path = ShowInput("Directory with install python not found, please enter manually")
	} else if countV == 1 {
		path = baseDir + versions[0]
		fmt.Printf("Find directory: %s\nUsing this directory by default\n.", path)
	} else {
		fmt.Println("Find more then one directory with installed python")
		path = baseDir + ShowChoices(versions)
	}

	return path
}

func getPythonDirectory(baseDir, prefix string) string {
	var versions []string
	dirInfo, err := ioutil.ReadDir(baseDir)
	if err != nil {
		panic(err)
	}

	for _, dir := range dirInfo {
		name := dir.Name()
		if strings.HasPrefix(name, prefix) && dir.IsDir() {
			versions = append(versions, name)
		}
	}

	return SelectOneDirectory(baseDir, versions)
}

// function checks if directory "site-packages" exists, and if it does not, it creates
func createSitePackagesDir(path string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			color.PrintfWarning("site-packages directory not detected. Creating..\n")
			if err = os.MkdirAll(path, cfg.DEFAULT_CHMOD); err != nil {
				color.PrintfError("Can't create \"site-packages\" directory. Please create it manually, otherwise you will get an error\n")
			}
		} else {
			color.PrintfError("Can't check exists of the \"site-packages\" directory: %s\n", err.Error())
		}
	}
}

func GetPythonBinDirectory(version string) string {
	// TODO add venv support
	if OS == "windows" {
		// TODO
	} else if OS == "linux" {
		return "/usr/bin/python" + version
	} else if OS == "" {
		// TODO
	}

	panic("Unsupported OS\n")
}

// select one dir without version
func GetPythonLibDirectory() string {
	var path string
	if OS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		path = getPythonDirectory(homeDir + "\\AppData\\Local\\Programs\\Python\\", "Python") + "\\Lib\\site-packages"
	} else if OS == "linux" {
		path = getPythonDirectory("/usr/lib/", "python") + "/site-packages"
	} else if OS == "darwin" { // mac os
		path = getPythonDirectory("/usr/local/lib/", "python") + "/site-packages"
	} else {
		panic("Unsupported OS\n")
	}

	createSitePackagesDir(path)
	return path
}

func ParsePythonVersion(path string) string {
	re, _ := regexp.Compile(`\d[\d|\.]+`)
	v := re.FindString(path)

	if OS == "windows" {
		if v == "" {
			panic("Unable to parse version. Please check, that the selected directory contains version")
		}

		v = fmt.Sprint(v[0], ".", v[1:])
	}

	return v
}

// format name for correct search
func FormatName(dirName string) string {
	return strings.ToLower(strings.ReplaceAll(dirName, "-", "_"))
}

func GetFirstPackageMetaDir(pkgName string) string {
	dirInfo, err := ioutil.ReadDir(cfg.PathToLib)
	if err != nil {
		return ""
	}
	pkgName = FormatName(pkgName)

	for _, dir := range dirInfo {
		curPkgName, v, _ := ParseDirectoryName(dir.Name())
		if FormatName(curPkgName) == pkgName && v != "" {
			return dir.Name()
		}
	}

	return ""
}

func GetAllPackageDirs(pkgName string) []string {
	var dirs []string
	dirInfo, err := ioutil.ReadDir(cfg.PathToLib)
	if err != nil {
		return dirs
	}
	pkgName = FormatName(pkgName)

	// first we check if we have found the right name, then we check if we have exceeded the boundaries
	var originalName string
	for _, dir := range dirInfo {
		originalName = dir.Name()
		curPkgName, _, _ := ParseDirectoryName(originalName)
		if FormatName(curPkgName) == pkgName {
			// TODO optimize this shit
			dirs = append(dirs, originalName)
		}
	}

	return dirs
}

// Count size of all files in directory and return bytes
func GetDirSize(dir string) int64 {
	var size int64
	_ = filepath.Walk(cfg.PathToLib + dir, func(_ string, info os.FileInfo, _ error) error {
		if info != nil && !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	return size
}

