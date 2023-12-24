package web

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
	"github.com/fextpkg/cli/fext/ferror"
)

type PyPiRequest struct {
	// Input package name
	pkgName string
	// Input conditions
	conditions []expression.Condition
}

// packageTags stores package compatibility tags
// (https://peps.python.org/pep-0425/)
type packageTags struct {
	// Package name
	name string
	// Package version
	version string
	// Optional tag that is rarely used
	buildTag string
	// The Python tag indicates the implementation and version required by a
	// distribution
	pyTag string
	// The ABI tag indicates which Python ABI is required by any included
	// extension modules. For implementation-specific ABIs, the implementation
	// is abbreviated in the same way as the Python Tag, e.g., cp33d would be
	// the CPython 3.3 ABI with debugging
	abiTag string
	// The platform tag is simply distutils.util.get_platform() with all
	// hyphens "-" and periods ".", replaced with underscore "_"
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

// DownloadPackage downloads the package file from PyPi repository.
// Returns a path to downloaded file
func (req *PyPiRequest) DownloadPackage(link string) (string, error) {
	hashSum := strings.Split(link, "sha256=")[1]

	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.Create(filepath.Join(config.PythonLibPath, hashSum+".tmp"))
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err = io.Copy(tmpFile, io.Reader(resp.Body)); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// getPackageList gets a web page with the package list
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

// Parse document and select a correct version. Returns version, download link
func (req *PyPiRequest) selectSuitableVersion(doc *html.Node) (string, string, error) {
	// html => body (on pypi)
	startNode := doc.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.LastChild

	// Check the latest versions first
	for node := startNode; node != nil; node = node.PrevSibling {
		// Elements with package data are stored in the "a" tag.
		// Exclude all other tags, like "br" and others.
		if node.Data != "a" {
			continue
		}

		version, link, err := req.getPackageInfo(node)
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

// getPackageInfo parses the node and checks all the data from it for
// compliance with the desired version. If successful, it returns the version and
// downloads link. If the version could not be found, empty strings will be
// returned without an error. If an error occurred, it will be returned with
// empty strings.
func (req *PyPiRequest) getPackageInfo(node *html.Node) (string, string, error) {
	fullData := node.FirstChild.Data

	// Select only wheel package
	if !strings.HasSuffix(fullData, ".whl") {
		return "", "", nil
	}

	pkgTags := parsePackageTags(fullData)

	// Check Python compatibility tags
	ok, err := pkgTags.CheckCompatibility()
	if !ok {
		return "", "", err
	}

	// Check package version
	ok, err = expression.CompareConditions(pkgTags.version, req.conditions)
	if !ok {
		if errors.Is(err, strconv.ErrSyntax) {
			// Due to the very strange description and lack of compatibility
			// with semantic versioning in PEP 440, we don't have an elegant
			// and efficient algorithm for handling post and dev releases.
			// Therefore, if we encounter an error, it is most likely related
			// to parsing these two parts. Since it is undesirable to interrupt
			// the package downloading, we ignore this case. These two
			// additional parts of the version are not highly significant.
			return "", "", nil
		}
		return "", "", err
	}

	link, versionRequirements := parseAttrs(node.Attr)
	_, conditions := expression.ParseConditions(versionRequirements)

	// Check the Python version
	ok, err = expression.CompareConditions(config.PythonVersion, conditions)
	if !ok {
		return "", "", err
	}

	return pkgTags.version, link, nil
}

// checkPythonCompatibility accepts python-tag of a package (PEP 425) and checks
// compatibility with installed python. Currently, only "py" and "cp" tags are
// implemented. When any other tag is passed, false will be returned.
func (tag *packageTags) checkPythonCompatibility() bool {
	// https://packaging.python.org/en/latest/specifications/platform-compatibility-tags/#python-tag

	// There can be several versions and they alternate through a point
	for _, version := range strings.Split(tag.pyTag, ".") {
		code := version[:2]

		if code == "py" {
			// Since there is support for versions 2 and 3 at once, we only check for the
			// presence of the number 3
			if version[2:] == "3" {
				return true
			}
		} else if code == "cp" {
			cpythonVersion := config.GetPythonMinorVersion()
			// Remove the extra characters and compare only the minor version
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
func (tag *packageTags) checkPlatformCompatibility() (bool, error) {
	// https://packaging.python.org/en/latest/specifications/platform-compatibility-tags/#platform-tag
	if tag.platformTag == "any" {
		return true, nil
	}
	return checkCompatibility(tag.platformTag)
}

// CheckCompatibility calls other tag checking methods and returns the final
// result.
func (tag *packageTags) CheckCompatibility() (bool, error) {
	ok, err := tag.checkPlatformCompatibility()
	if err != nil || !ok {
		return false, err
	}

	return tag.checkPythonCompatibility(), nil
}

// NewRequest creates a new package search query object on PyPi with the
// specified conditions
func NewRequest(pkgName string, cond []expression.Condition) *PyPiRequest {
	return &PyPiRequest{
		pkgName:    pkgName,
		conditions: cond,
	}
}

// parsePackageTags separates all package tags and creates a new structure
// packageTags with them.
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

// parseAttrs parses the HTML element attributes and return download link,
// python requirement versions. Example: ("https://...", ">=3.7")
func parseAttrs(attrs []html.Attribute) (string, string) {
	var link, versionRequirements string
	for _, attr := range attrs {
		switch attr.Key {
		case "href":
			link = attr.Val
		case "data-requires-python":
			versionRequirements = strings.ReplaceAll(attr.Val, ",", "")
		}
	}
	return link, versionRequirements
}
