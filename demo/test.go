package main

import (
	"fmt"
	"runtime"
)

type User struct {
	Name string
}

func setUser(list []*User) {
	for _, v := range list {
		v.Name = "lisi"
	}
}

func main() {
	//runtime.GOARCH 返回当前的系统架构；runtime.GOOS 返回当前的操作系统。
	sysType := runtime.GOOS

	if sysType == "linux" {
		// LINUX系统
		fmt.Println("Linux system")
	}

	if sysType == "windows" {
		// windows系统
		fmt.Println("Windows system")
	}
}
