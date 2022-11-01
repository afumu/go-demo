package main

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

type User struct {
	Name string
}

func main() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("发生错误了")
				buf := new(bytes.Buffer)
				fmt.Fprintf(buf, "%v\n", r)
				for i := 1; ; i++ {
					pc, file, line, ok := runtime.Caller(i)
					f := runtime.FuncForPC(pc)
					if !ok {
						break
					}
					fmt.Fprintf(buf, "%s:%d (0x%x) Cause by: %s\n", file, line, pc, f.Name())
				}
				fmt.Println(buf.String())

				/*			_, fileName, lineNumber, ok := runtime.Caller(4)
							if ok {
								lineNumberStr := fmt.Sprintf("%d", lineNumber)
								message := fileName + ":" + lineNumberStr + " "
								errMessage := fmt.Sprintf("%v", r)
								message += errMessage
								fmt.Println(message)
							}*/
			}
		}()
		fmt.Println("111111")
		fmt.Println("1111")
		fmt.Println("111")
		fmt.Println("1111111111")
		time.Sleep(2 * time.Second)
		fmt.Println("----------------------")
		fmt.Println("----------------------")
		seek()
	}()

	select {}
}

func seek() {
	fmt.Println("----------------------")
	fmt.Println("----------------------")
	var user *User
	//panic("aaaaaaaaaaaaaa")
	name := user.Name
	fmt.Println(name)
}
