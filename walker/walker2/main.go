package main

import (
	"bufio"
	"errors"
	"fmt"
	dio "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/saracen/walker"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	//var characters = "机密,绝密"
	//search, data := LinuxFileSearch(characters)
	//fmt.Println(search)
	//fmt.Println(data)
	//if len(os.Args) == 2 {
	//	path := os.Args[1]
	//	read(path)
	//	return
	//}
	//
	path2 := os.Args[1]
	LinuxFileSearch(path2, "机密,绝密")
	//read("")

	//home/client/秘密aa.txt
	//data/home/client/秘密aa.txt
	//opt/机密aa.txt
	///home/client/秘密aa.txt

	/*	var a = "//home/client/秘密aa.txt"
		all := strings.ReplaceAll(a, "//", "/")
		fmt.Println(all)*/

}

/*
func read(dir string) {
	/*	dirs, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		for i, v := range dirs {
			if v.IsDir() {
				join := filepath.Join(dir, v.Name())
				read(join)
			}
			join := filepath.Join(dir, v.Name())
			message := fmt.Sprintf("%d-%v", i, join)
			fmt.Println(message)
		}*/
/*	var path = "E:\\vuln"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}*/
//}

func LinuxFileSearch(path, characters string) (string, string) {
	var isPass, data = "true", ""
	var mu sync.Mutex
	split := strings.Split(path, ",")
	for _, p := range split {
		if len(strings.TrimSpace(characters)) > 0 {
			list := strings.Split(strings.TrimSpace(characters), ",")
			walkFn := func(pathname string, fi os.FileInfo) error {
				opener := fileOpener(pathname)

				if strings.Contains(pathname, "proc") || strings.Contains(pathname, "sys") ||
					strings.Contains(pathname, "dev") || strings.Contains(pathname, "run") {
					return filepath.SkipDir
				}

				// 打开文件（文件路径外部函数已经赋值好了）
				rc, err := opener()

				if errors.Is(err, fs.ErrPermission) {
					// 如果发生读取文件没有权限的话，直接结束for循环
					log.Println("Permission error:", pathname)
					return nil
				} else if err != nil {
					log.Println("unable to open", pathname, err)
					return nil
				}

				go func() {
					defer rc.Close()
					if !fi.IsDir() {
						line := bufio.NewReader(rc)
						_, _, _ = line.ReadLine()
						for _, character := range list {
							if strings.Contains(fi.Name(), character) {
								mu.Lock()
								isPass = "false"
								data += fmt.Sprintf("%v\n", pathname)
								log.Println(pathname)
								mu.Unlock()
							}
						}
					}
				}()

				return nil
			}
			errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
				if os.IsPermission(err) {
					return nil
				}
				return err
			})
			err := walker.Walk(p, walkFn, errorCallbackOption)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return isPass, data
}

var ch = make(chan struct{}, 4)

func readAllDir(path string) {
	ch <- struct{}{}
	defer func() {
		<-ch
	}()
	dirs, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	for i, v := range dirs {
		if v.IsDir() {
			join := filepath.Join(path, v.Name())
			readAllDir(join)
		}
		join := filepath.Join(path, v.Name())
		message := fmt.Sprintf("%d-%v", i, join)
		fmt.Println(message)
	}
}

func fileOpener(pathname string) func() (dio.ReadSeekCloserAt, error) {
	return func() (dio.ReadSeekCloserAt, error) {
		return os.Open(pathname)
	}
}
