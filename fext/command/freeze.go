package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fextpkg/cli/fext/config"
	"github.com/fextpkg/cli/fext/pkg"
)

func Freeze() {
	files, err := os.ReadDir(config.PythonLibPath)

	if err != nil {
		log.Fatal(err)
	}

	var count int
	var size int64
	for _, f := range files {
		dirName := f.Name()
		if f.IsDir() && strings.HasSuffix(dirName, "dist-info") {
			p, err := pkg.LoadFromMetaDir(dirName)
			if err == nil {
				count++
				s, _ := p.GetSize()
				size += s
				fmt.Printf("%s (%s)\n", p.Name, p.Version)
			}
		}
	}

	fmt.Printf("\nTotal: %d (%.2f MB)\n", count, float32(size/1024)/1024)
}
