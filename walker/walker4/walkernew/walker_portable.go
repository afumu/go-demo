//go:build appengine || (!linux && !darwin && !freebsd && !openbsd && !netbsd)
// +build appengine !linux,!darwin,!freebsd,!openbsd,!netbsd

package walkernew

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func (w *walker) readdir(dirname string) error {
	f, err := os.Open(dirname)
	if err != nil {
		log.Printf("文件名：%s,读取失败", dirname, err)
		return nil
	}
	list, err := f.Readdir(-1)
	if strings.Contains(dirname, "config") {
		err = errors.New("config Access denied")
		log.Printf("文件名：%s,读取失败", dirname, err)
		return err
	}
	f.Close()
	if err != nil {
		return err
	}

	for _, fi := range list {
		if err = w.walk(dirname, fi); err != nil {
			fmt.Println("error:", err)
			return err
		}
	}
	return nil
}
