package pkg

import (
	"github.com/fextpkg/cli/fext/config"

	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func unzip(path string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()
	path = filepath.Dir(path) // remove hashsum part

	for _, f := range r.File {
		fpath := filepath.Join(path, f.Name)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, config.DefaultChmod)
			if err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), config.DefaultChmod); err != nil {
				return err
			} else {
				outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, config.DefaultChmod)
				if err != nil {
					return err
				}

				rf, err := f.Open()
				if err != nil {
					return err
				}

				_, err = io.Copy(outFile, rf)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Extract(path string) error {
	return unzip(path)
}
