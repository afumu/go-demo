package main

import (
	"fmt"
	dio "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/java/jar"
	"github.com/saracen/walker"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	path := os.Args[1]
	fmt.Println(path)
	walkFn := func(pathname string, fi os.FileInfo) error {
		opener := fileOpener(pathname)

		if strings.Contains(pathname, "proc") ||
			strings.Contains(pathname, "dev") || strings.Contains(pathname, "run") {
			return filepath.SkipDir
		}

		// 打开文件（文件路径外部函数已经赋值好了）
		rc, err := opener()
		if err != nil {
			fmt.Println("打开文件错误：", err)
		}

		p := jar.NewParser(jar.WithSize(fi.Size()), jar.WithFilePath(pathname), jar.WithOffline(true))

		libs, deps, err := p.Parse(rc)
		fmt.Printf("%+v\n", libs)
		fmt.Printf("%+v\n", deps)
		return nil
	}
	errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
		if os.IsPermission(err) {
			return nil
		}
		return err
	})
	err := walker.Walk(path, walkFn, errorCallbackOption)
	if err != nil {
		fmt.Println(err)
	}
}

func fileOpener(pathname string) func() (dio.ReadSeekCloserAt, error) {
	return func() (dio.ReadSeekCloserAt, error) {
		return os.OpenFile(pathname, syscall.O_DIRECT, 0666)
	}
}
