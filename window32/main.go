package main

import (
	"fmt"
	"github.com/StackExchange/wmi"
)

type DiskInfo struct {
	// 全部字段请看文档,这里我只取了硬盘名字,大小,序列号
	// https://docs.microsoft.com/zh-cn/windows/win32/cimwin32prov/win32-diskdrive#syntax
	Name         string
	Size         uint64
	SerialNumber string
	SystemName   string
}

func getDiskInfo() {
	var diskInfos []DiskInfo
	err := wmi.Query("Select * from Win32_DiskDrive", &diskInfos)
	if err != nil {
		return
	}
	for _, disk := range diskInfos {
		fmt.Printf("硬盘名称是%s,硬盘大小是%dG,硬盘序列号是%s\n", disk.Name, disk.Size/1024/1024/1024, disk.SerialNumber)
	}
}
func main() {
	// TODO 先下载wmi库
	// go get github.com/StackExchange/wmi
	getDiskInfo()
}
