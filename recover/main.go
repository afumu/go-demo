package main

import (
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
				fmt.Println(r)
				caller, file, line, ok := runtime.Caller(0)
				fmt.Println(caller, file, line, ok)
			}
		}()
		fmt.Println("111111")
		fmt.Println("1111")
		fmt.Println("111")
		fmt.Println("1111111111")
		time.Sleep(5 * time.Second)
		fmt.Println("----------------------")
		fmt.Println("----------------------")
		fmt.Println("----------------------")
		fmt.Println("----------------------")
		var user *User
		name := user.Name
		fmt.Println(name)
	}()

	select {}
}
