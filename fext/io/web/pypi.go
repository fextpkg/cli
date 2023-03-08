package web

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/expression"
)

type PyPi struct {
	pkgName string
	op      [][]string
}

func (web *PyPi) GetPackageData() (string, string, error) {
	doc, err := web.getPackageList()
	if err != nil {
		return "", "", err
	}

	return web.selectAppropriateVersion(doc)
}

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

// Parse document and select optimal version. Returns version and link to download
func (web *PyPi) selectAppropriateVersion(doc *html.Node) (string, string, error) {
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

		pkgData := strings.Split(fullData, "-") // [name, version, [,build-tag] py-tag, abi-tag, platform-tag]
		ok, err := checkPlatformCompatibility(pkgData)
		if !ok {
			continue
		}

		ok, err = compareVersion(pkgData[1], web.op)
		if !ok {
			if err != nil {
				return "", "", err
			}
			continue
		}

		link, versionClassifiers := parseAttrs(node.Attr)
		_, classifiers := expression.ParseExpression(versionClassifiers)
		ok, err = compareVersion(config.PythonVersion, classifiers)
		if !ok {
			if err != nil {
				return "", "", err
			}
			continue
		}

		return pkgData[1], link, nil
	}

	return "", "", errors.New("no matching version was found")
}

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

func NewRequest(pkgName string, op [][]string) *PyPi {
	return &PyPi{
		pkgName: pkgName,
		op:      op,
	}
}

func parseAttrs(attrs []html.Attribute) (string, string) {
	var link, versionClassifiers string
	for _, attr := range attrs {
		switch attr.Key {
		case "href":
			link = attr.Val
		case "data-requires-python":
			// remove this parts, because it's works fine without it
			attr.Val = strings.ReplaceAll(attr.Val, ".*", "")
			versionClassifiers = attr.Val
		}
	}
	return link, versionClassifiers
}

func compareVersion(version string, operators [][]string) (bool, error) {
	for _, op := range operators {
		ok, err := expression.CompareVersion(version, op[0], op[1])
		if !ok {
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}
	return true, nil
}

func checkPlatformCompatibility(pkgData []string) (bool, error) {
	var platformTag string
	if len(pkgData) == 6 { // have additional build tag
		platformTag = pkgData[5]
	} else {
		platformTag = pkgData[4]
	}
	platformTag = platformTag[:len(platformTag)-4] // remove ".whl"
	if platformTag == "any" {
		return true, nil
	}
	return checkCompatibility(platformTag)
}
