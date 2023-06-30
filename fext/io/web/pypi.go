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

type PyPi struct {
	pkgName    string
	conditions []expression.Condition
}

// GetPackageData gets a first package version that fits the conditions of the
// operators. Returns package version, download link. An error will be returned
// if a suitable version was not found or another error occurred
func (web *PyPi) GetPackageData() (string, string, error) {
	doc, err := web.getPackageList()
	if err != nil {
		return "", "", err
	}

	return web.selectSuitableVersion(doc)
}

// getPackageList gets a list of package versions
func (web *PyPi) getPackageList() (*html.Node, error) {
	resp, err := http.Get("https://pypi.org/simple/" + web.pkgName + "/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		return nil, err
	}

	return doc, nil
}

// Parse document and select correct version. Returns version, download link
func (web *PyPi) selectSuitableVersion(doc *html.Node) (string, string, error) {
	// html => body (on pypi)
	startNode := doc.FirstChild.NextSibling.FirstChild.NextSibling.NextSibling.LastChild
	var fullData string

	// check latest versions first
	for node := startNode; node != nil; node = node.PrevSibling {
		if node.Data != "a" {
			continue
		} else {
			fullData = node.FirstChild.Data
		}
		// select only wheel
		if !strings.HasSuffix(fullData, ".whl") {
			continue
		}

		pkgTags := strings.Split(fullData, "-") // [name, version, [,build-tag] py-tag, abi-tag, platform-tag]
		ok, err := checkPlatformCompatibility(pkgTags)
		if !ok {
			if err != nil {
				return "", "", err
			}
			continue
		}

		ok, err = compareVersion(pkgTags[1], web.conditions)
		if !ok {
			if err != nil {
				return "", "", err
			}
			continue
		}

		link, versionClassifiers := parseAttrs(node.Attr)
		_, classifiers := expression.ParseConditions(versionClassifiers)
		ok, err = compareVersion(config.PythonVersion, classifiers)
		if !ok {
			if err != nil {
				return "", "", err
			}
			continue
		}

		return pkgTags[1], link, nil
	}

	return "", "", ferror.NoSuitableVersion
}

// DownloadPackage downloads the package from PyPi repository. Returns path
// to downloaded package
func (web *PyPi) DownloadPackage(link string) (string, error) {
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
		tmpFile.Close()
		return "", err
	}

	return tmpFile.Name(), nil
}

// NewRequest creates a new package search query object on PyPi with the
// specified conditions
func NewRequest(pkgName string, cond []expression.Condition) *PyPi {
	return &PyPi{
		pkgName:    pkgName,
		conditions: cond,
	}
}

// parseAttrs parses the HTML element attributes and returns download link,
// version classifiers
func parseAttrs(attrs []html.Attribute) (string, string) {
	var link, versionClassifiers string
	for _, attr := range attrs {
		switch attr.Key {
		case "href":
			link = attr.Val
		case "data-requires-python":
			// remove this parts, because it's works fine without it
			versionClassifiers = strings.ReplaceAll(attr.Val, ".*", "")
		}
	}
	return link, versionClassifiers
}

// compareVersion checks the compliance of the version for the passed operators
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

// checkPlatformCompatibility accepts compatibility tags of a package (PEP 425) and checks
// the platform tag for compatibility with the current platform
func checkPlatformCompatibility(tags []string) (bool, error) {
	var platformTag string
	if len(tags) == 6 { // have additional build tag
		platformTag = tags[5]
	} else {
		platformTag = tags[4]
	}
	platformTag = platformTag[:len(platformTag)-4] // remove ".whl"
	if platformTag == "any" {
		return true, nil
	}
	return checkCompatibility(platformTag)
}
