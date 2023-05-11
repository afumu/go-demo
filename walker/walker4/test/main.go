package main

import (
	"fmt"
	"log"
	"os"
	"test-demo/walker/walker4/walkernew"
)

func main() {
	path := "F:\\agent"
	walkFn := func(pathname string, fi os.FileInfo) error {
		log.Printf("%v\n", pathname)
		//opener := fileOpener(pathname)

		// 打开文件（文件路径外部函数已经赋值好了）
		//rc, err := opener()
		//if err != nil {
		//	fmt.Println("打开文件错误：", err)
		//}

		//p := jar.NewParser(jar.WithSize(fi.Size()), jar.WithFilePath(pathname), jar.WithOffline(true))
		//libs, deps, err := p.Parse(rc)
		//fmt.Printf("%+v\n", libs)
		//fmt.Printf("%+v\n", deps)
		return nil
	}
	errorCallbackOption := walkernew.WithErrorCallback(func(pathname string, err error) error {
		if os.IsPermission(err) {
			return nil
		}
		return err
	})
	err := walkernew.Walk(path, walkFn, errorCallbackOption)
	if err != nil {
		fmt.Println(err)
	}
}
