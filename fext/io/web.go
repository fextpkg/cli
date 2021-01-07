package io

import (
	"github.com/Flacy/fext/fext/cfg"
	"github.com/Flacy/fext/fext/utils"

	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getPackageList(name string) (*html.Node, error) {
	resp, err := http.Get(cfg.BASE_PACKAGE_URL + name + "/")
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

func downloadPackage(buffer interface {
	Write([]byte) (int, error)
	UpdateTotal(int)
},
	link string) (string, error) {
	hashSum := strings.Split(link, "sha256=")[1]

	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.Create(cfg.PathToLib + hashSum + ".tmp")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	buffer.UpdateTotal(int(resp.ContentLength))

	if _, err = io.Copy(tmpFile, io.TeeReader(resp.Body, buffer)); err != nil {
		tmpFile.Close()
		return "", err
	}

	return tmpFile.Name(), nil
}

// Parse document and select optimal version. Returns package version, link to download
func selectCorrectPackageVersion(doc *html.Node, op [][]string) (string, string, error) {
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
		pkgVersion := strings.Split(fullData, "-")[1] // [name, version, ...]

		var link string
		var versionClassifiers []string
		for _, attr := range node.Attr {
			switch attr.Key {
			case "href":
				link = attr.Val
			case "data-requires-python":
				// remove this parts, cause it's works fine without it
				attr.Val = strings.ReplaceAll(attr.Val, ".*", "")
				versionClassifiers = strings.Split(attr.Val, ", ")
			}
		}

		// if user specified compare operator
		hasNoOk := false // flag used for continue parent loop
		for _, op := range op {
			ok, err := utils.CompareVersion(pkgVersion, op[0], op[1]) // [op, version]
			if err != nil {
				return "", "", err
			} else if !ok {
				hasNoOk = true
				break
			}
		}
		if hasNoOk {
			continue
		}

		for _, classifier := range versionClassifiers {
			_, op := utils.SplitOperators(classifier)
			if ok, err := utils.CompareVersion(cfg.PythonVersion, op[0][0], op[0][1]); !ok {
				continue
			} else if err != nil {
				return "", "", err
			}
		}
		return pkgVersion, link, nil
	}

	return "", "", errors.New("No matching version was found")
}

