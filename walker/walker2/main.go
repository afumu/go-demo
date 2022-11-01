package main

import (
	"fmt"
	"github.com/saracen/walker"
	"log"
	"os"
	"strings"
)

func main() {
	var characters = "机密,绝密"
	search, data := LinuxFileSearch(characters)
	fmt.Println(search)
	fmt.Println(data)
}

func LinuxFileSearch(characters string) (string, string) {
	var isPass, data = "true", ""
	if len(strings.TrimSpace(characters)) > 0 {
		list := strings.Split(strings.TrimSpace(characters), ",")
		walkFn := func(pathname string, fi os.FileInfo) error {
			if !fi.IsDir() {
				for _, character := range list {
					if strings.Contains(fi.Name(), character) {
						isPass = "false"
						data += fmt.Sprintf("%v\n", pathname)
						log.Println(pathname)
					}
				}
			}
			return nil
		}
		errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
			if os.IsPermission(err) {
				return nil
			}
			return err
		})
		_ = walker.Walk("/", walkFn, errorCallbackOption)
	}
	return isPass, data
}
