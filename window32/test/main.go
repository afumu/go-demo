package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	number := getSerialNumber()
	fmt.Println(number)

	/*	//result := executeCommand("lsblk","","")
		result := executeCommand("lsblk")
		//result := executeCommand("lsblk", "| awk {print $1}")
		fmt.Println(result)
		arr := strings.Split(result, "\n")
		for _, v := range arr {
			if strings.HasPrefix(v, "s") {
				split := strings.Split(v, " ")
				fmt.Println(split[0])
			}
		}*/

	/*	var str = "sda\n├─sda1\n├─sda2\n└─sda3\n├─zstack-root\n└─zstack-swap\n"

		arr := strings.Split(str, "\n")
		for _, v := range arr {
			if strings.HasPrefix(v, "s") {
				fmt.Println(v)
			}
		}*/

}

func getSerialNumber() string {
	// 获取硬盘序列号
	var serialNumber string
	var linuxSerialNumbers []string
	if runtime.GOOS == "windows" {
		result := executeCommand("wmic", "diskdrive", "get", "SerialNumber")
		serialNumber = strings.ReplaceAll(result, "SerialNumber", "")
	} else {
		// lsblk |awk '{print $1}'
		result := executeCommand("lsblk")
		arr := strings.Split(result, "\n")
		for _, v := range arr {
			if strings.HasPrefix(v, "s") {
				// lsblk -n --nodeps -o name,serial /dev/sda
				split := strings.Split(v, " ")
				if len(split) > 0 {
					sn := executeCommand("lsblk", "-n", "--nodeps", "-o", "name,serial", "/dev/"+split[0])
					linuxSerialNumbers = append(linuxSerialNumbers, sn)
				}
			}
		}
		serialNumber = strings.Join(linuxSerialNumbers, ",")
	}
	return serialNumber
}

func executeCommand(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	return out.String()
}
