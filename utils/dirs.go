package utils

import (
	"fmt"
	"io/ioutil"
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
	dirInfo, _ := ioutil.ReadDir(baseDir)

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
		return getPythonDirectory("C:\\", "Python") + "\\Lib\\site-packages"
	} else if OS == "linux" {
		return getPythonDirectory("/usr/lib/", "python") + "/site-packages"
	} else if OS == "" {
		// TODO
	}

	panic("Unsupported OS\n")
}

func ParsePythonVersion(path string) string {
	re, _ := regexp.Compile(`\d[\d|\.]+`)
	v := re.FindString(path)
	if OS == "windows" {
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
		} else {
			findPrefix = true
			dirs = append(dirs, dirName)
		}
	}
	
	return dirs
}

