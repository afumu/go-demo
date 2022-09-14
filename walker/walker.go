package main

import (
	"fmt"
	"github.com/aquasecurity/go-dep-parser/pkg/java/jar"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/saracen/walker"
	"os"
	"path/filepath"
	"strings"
)

var requiredExtensions = []string{".jar", ".war", ".ear", ".par"}

func main() {
	var libs []types.Library
	var deps []types.Dependency
	walkFn := func(pathname string, fi os.FileInfo) error {
		fmt.Printf("%s: %d bytes\n", pathname, fi.Size())
		ext := filepath.Ext(pathname)
		for _, required := range requiredExtensions {
			if strings.EqualFold(ext, required) {
				file, err := os.Open(pathname)
				if err != nil {
					fmt.Println(err.Error())
				}
				fileInfo, err := file.Stat()
				if err != nil {
					fmt.Println(err.Error())
				}
				p := jar.NewParser(jar.WithSize(fileInfo.Size()), jar.WithFilePath(pathname), jar.WithOffline(true))

				libs, deps, err = p.Parse(file)
			}
		}
		return nil
	}

	// error function called for every error encountered
	errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
		if os.IsPermission(err) {
			return nil
		}
		return err
	})

	walker.Walk("D:\\workplace\\junan-template", walkFn, errorCallbackOption)

	fmt.Printf("%+v", libs)
	fmt.Printf("%+v", deps)
}
