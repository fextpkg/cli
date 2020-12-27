package utils

import (
	"fmt"
	"io/ioutil"
	"os"
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
func FindPythonLibDirectory() string {
	if OS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		return getPythonDirectory(homeDir + "\\AppData\\Local\\Programs\\Python\\", "Python") + "\\Lib\\site-packages"
	} else if OS == "linux" {
		return getPythonDirectory("/usr/lib/", "python") + "/site-packages"
	} else if OS == "darwin" { // mac os
		return getPythonDirectory("/usr/local/lib/", "python") + "/site-packages"
	}

	panic("Unsupported OS\n")
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

func GetFirstPackageMetaDir(libDir, pkgName string) string {
	dirInfo, err := ioutil.ReadDir(libDir)
	if err != nil {
		return ""
	}
	pkgName = strings.ReplaceAll(pkgName, "-", "_")

	for _, dir := range dirInfo {
		dirName := dir.Name()
		if strings.HasPrefix(dirName, pkgName + "-") {
			return dirName
		}
	}

	return ""
}

func GetAllPackageDirs(pkgName, libDir string) []string {
	var dirs []string
	dirInfo, err := ioutil.ReadDir(libDir)
	if err != nil {
		return dirs
	}
	pkgName = strings.ReplaceAll(pkgName, "-", "_")

	// first we check if we have found the right name, then we check if we have exceeded the boundaries
	var findPrefix bool
	for _, dir := range dirInfo {
		dirName := dir.Name()
		if !strings.HasPrefix(dirName, pkgName) {
			if findPrefix {
				break
			}
		} else if dir.IsDir() {
			findPrefix = true
			dirs = append(dirs, dirName)
		}
	}

	return dirs
}

