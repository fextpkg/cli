package web

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fextpkg/cli/fext/ferror"
	"golang.org/x/net/html"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
)

type PyPiRequest struct {
	pkgName    string
	conditions []expression.Condition
}

type packageTags struct {
	name        string
	version     string
	buildTag    string
	pyTag       string
	abiTag      string
	platformTag string
}

// GetPackageData gets a first package version that fits the conditions of the
// operators and system requirements. Returns package version, download link. An
// error will be returned if a suitable version was not found or another error
// occurred
func (req *PyPiRequest) GetPackageData() (string, string, error) {
	doc, err := req.getPackageList()
	if err != nil {
		return "", "", err
	}

	return req.selectSuitableVersion(doc)
}

// DownloadPackage downloads the package from PyPi repository. Returns path to
// downloaded package
func (req *PyPiRequest) DownloadPackage(link string) (string, error) {
	hashSum := strings.Split(link, "sha256=")[1]

	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.Create(config.PythonLibPath + hashSum + ".tmp")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, io.Reader(resp.Body)); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// getPackageList gets a web page with package list
func (req *PyPiRequest) getPackageList() (*html.Node, error) {
	resp, err := http.Get("https://pypi.org/simple/" + req.pkgName + "/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(strings.ToLower(resp.Status[4:]))
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// Parse document and select correct version. Returns version, download link
func (req *PyPiRequest) selectSuitableVersion(doc *html.Node) (string, string, error) {
	// html => body (on pypi)
	startNode := doc.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.LastChild

	// check latest versions first
	for node := startNode; node != nil; node = node.PrevSibling {
		// Elements with package data are stored in "a" attribute. Exclude all other attrs like "br" and others.
		if node.Data != "a" {
			continue
		}

		version, link, err := req.checkPackageInfo(node)
		if err != nil {
			// Critical error, it is impossible to continue the search
			return "", "", err
		} else if version == "" {
			// A suitable version was not found, continue the search
			continue
		} else {
			// A suitable version is found
			return version, link, nil
		}
	}

	return "", "", ferror.NoSuitableVersion
}

// checkPackageInfo parses the node and checks all the data from it for
// compliance with the desired version. If successful, it returns the version and
// download link. If the version could not be found, empty strings will be
// returned without an error. If an error occurred, it will be returned with
// empty strings.
func (req *PyPiRequest) checkPackageInfo(node *html.Node) (string, string, error) {
	fullData := node.FirstChild.Data

	// Select only wheel package
	if !strings.HasSuffix(fullData, ".whl") {
		return "", "", nil
	}

	pkgTags := parsePackageTags(fullData)

	// Check platform-tag
	ok, err := checkPlatformCompatibility(pkgTags.platformTag)
	if !ok {
		return "", "", err
	}

	// Check python tag
	ok = checkPythonCompatibility(pkgTags.pyTag)
	if !ok {
		return "", "", nil
	}

	// Check package version
	ok, err = compareVersion(pkgTags.version, req.conditions)
	if !ok {
		return "", "", err
	}

	link, versionRequirements := parseAttrs(node.Attr)
	_, conditions := expression.ParseConditions(versionRequirements)

	// Check python version
	ok, err = compareVersion(config.PythonVersion, conditions)
	if !ok {
		return "", "", err
	}

	return pkgTags.version, link, nil
}

// NewRequest creates a new package search query object on PyPiRequest with the
// specified conditions
func NewRequest(pkgName string, cond []expression.Condition) *PyPiRequest {
	return &PyPiRequest{
		pkgName:    pkgName,
		conditions: cond,
	}
}

func parsePackageTags(s string) *packageTags {
	var buildTag string
	var buildTagIndex int
	tags := strings.Split(s, "-") // [name, version, [,build-tag] py-tag, abi-tag, platform-tag]

	if len(tags) == 6 { // have optional build-tag
		buildTagIndex = 2
		buildTag = tags[buildTagIndex]
	} else {
		buildTagIndex = 1
		buildTag = ""
	}

	pkgTags := &packageTags{
		name:        tags[0],
		version:     tags[1],
		buildTag:    buildTag,
		pyTag:       tags[buildTagIndex+1],
		abiTag:      tags[buildTagIndex+2],
		platformTag: tags[buildTagIndex+3][:len(tags[buildTagIndex+3])-4], // remove ".whl"
	}

	return pkgTags
}

// parseAttrs parses the HTML element attributes and returns download link,
// python requirement versions. Example: ("https://...", ">=3.7")
func parseAttrs(attrs []html.Attribute) (string, string) {
	var link, versionRequirements string
	for _, attr := range attrs {
		switch attr.Key {
		case "href":
			link = attr.Val
		case "data-requires-python":
			// remove this parts, because it's works fine without it
			versionRequirements = strings.ReplaceAll(attr.Val, ".*", "")
		}
	}
	return link, versionRequirements
}

// compareVersion checks the compliance of the version for the passed operators.
// If all conditions are true, true will be returned, otherwise false. The error
// is returned in case of an incorrect operator or version.
func compareVersion(version string, conditions []expression.Condition) (bool, error) {
	for _, cond := range conditions {
		ok, err := expression.CompareVersion(version, cond.Op, cond.Value)
		if !ok {
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}
	return true, nil
}

// checkPythonCompatibility accepts python-tag of a package (PEP 425) and checks
// compatibility with installed python
func checkPythonCompatibility(pythonTag string) bool {
	// https://packaging.python.org/en/latest/specifications/platform-compatibility-tags/#python-tag

	// There can be several platforms and they alternate through a point
	for _, version := range strings.Split(pythonTag, ".") {
		code := version[:2]

		if code == "py" {
			// Since there is support for versions 2 and 3 at once, we only check for the
			// presence of the number 3
			if version[2:] == "3" {
				return true
			}
		} else if code == "cp" {
			// Remove the extra characters and compare only the minor version
			cpythonVersion := config.GetPythonMinorVersion()
			tagVersion := version[3:]

			if tagVersion == "" || tagVersion == cpythonVersion {
				return true
			}
		}
	}

	return false
}

// checkPlatformCompatibility accepts platform-tag of a package (PEP 425) and checks
// them for compatibility with the current platform
func checkPlatformCompatibility(platformTag string) (bool, error) {
	// https://packaging.python.org/en/latest/specifications/platform-compatibility-tags/#platform-tag
	if platformTag == "any" {
		return true, nil
	}
	return checkCompatibility(platformTag)
}
