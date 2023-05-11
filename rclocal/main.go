package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var path string

func init() {
	flag.StringVar(&path, "path", "", "程序路径")
}

func main() {
	flag.Parse()
	fmt.Println(path)

	install(path)
}

func install(commandPath string) {
	rcLocalContent, err := readRcLocal(commandPath)
	fmt.Println("rc.local 待写入内容:", rcLocalContent)
	file, err := os.OpenFile("/etc/rc.local", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开或创建rc.local错误:", err)
		return
	}
	defer file.Close()
	n, err := file.WriteString(rcLocalContent)
	if err != nil {
		fmt.Println("写入失败:", err)
	}
	fmt.Println("写入长度：", n)

	// 更改文件权限
	if err := os.Chmod("/etc/rc.local", 0755); err != nil {
		fmt.Println(err)
		return
	}
	// 启动一个goroutine来执行rc.local脚本
	go func() {
		cmd := exec.Command("/etc/rc.local")
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println("程序已设置为开机自启动。")
}

func readRcLocal(commandPath string) (string, error) {
	var rcLocalContent string
	exists, _ := pathExists("/etc/rc.local")
	if exists {
		rcLocalContentByte, err := os.ReadFile("/etc/rc.local")
		if err != nil {
			return "", err
		}
		rcLocalContent = string(rcLocalContentByte)
	} else {
		rcLocalContent = `#!/bin/sh -e
#
# rc.local
#
# This script is executed at the end of each multiuser runlevel.
# Make sure that the script will "" on success or any other
# value on error.
#
# In order to enable or disable this script just change the execution
# bits.
#
# By default this script does nothing.

exit 0
`
	}
	rcLocalContentNew := strings.Replace(rcLocalContent, "exit 0", "", -1)
	fmt.Println("rcLocalContentNew:", rcLocalContentNew)
	rcLocalContentNew = rcLocalContentNew + "\n"
	rcLocalContentNew = rcLocalContentNew + commandPath + " &\n"
	rcLocalContentNew = rcLocalContentNew + "exit 0\n"
	return rcLocalContentNew, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	if err != nil {
		fmt.Println(err)
	}
	return false, err
}
