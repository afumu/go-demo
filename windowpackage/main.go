package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"sync"
)

type Software struct {
	DisplayIcon      string
	Name             string
	Version          string
	InstallLocation  string
	Publisher        string
	UninstallString  string
	InstallationTime string
}

var query64Path = "Software\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
var queryPath = "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall"

func main() {
	allSoftware := QueryWindowsSoftware()
	fmt.Println(allSoftware)
}

func QueryWindowsSoftware() []Software {
	var softwareDetailsSli []Software
	queryKey := func(w *sync.WaitGroup, res *[]Software) {
		defer w.Done()

		// 打开注册表文件
		kQuery64, err := registry.OpenKey(registry.LOCAL_MACHINE, query64Path, registry.READ)
		if err != nil {
			log.Println(err)
			return
		}

		// 读取注册表所有的key
		keyNames, err := kQuery64.ReadSubKeyNames(0)
		if err != nil {
			log.Println(err)
			return
		}

		// 查询出query64Path下面的程序详情，并且添加到SoftwareDetails
		softwareDetailsSli = getSoftwareDetails(keyNames, query64Path)
		*res = append(*res, softwareDetailsSli...)

		k, err1 := registry.OpenKey(registry.LOCAL_MACHINE, queryPath, registry.READ)
		if err1 != nil {
			return
		}

		// 读取注册表所有的key
		keyNames, err1 = k.ReadSubKeyNames(0)
		if err1 != nil {
			return
		}

		*res = append(*res, getSoftwareDetails(keyNames, queryPath)...)
	}

	var res []Software
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	go queryKey(waitGroup, &res)
	waitGroup.Wait()
	return res
}

func getSoftwareDetails(appName []string, queryPath string) []Software {
	softwareDetails := []Software{}
	for _, value := range appName {
		kQuery64Details, err := registry.OpenKey(registry.LOCAL_MACHINE, queryPath+"\\"+value, registry.READ)
		if err != nil {
			continue
		}
		displayIcon, _, err := kQuery64Details.GetStringValue("DisplayIcon")
		displayName, v, err := kQuery64Details.GetStringValue("DisplayName")
		displayVersion, _, err := kQuery64Details.GetStringValue("DisplayVersion")
		installLocation, _, err := kQuery64Details.GetStringValue("InstallLocation")
		publisher, _, err := kQuery64Details.GetStringValue("Publisher")
		uninstallString, _, err := kQuery64Details.GetStringValue("UninstallString")
		installDate, _, err := kQuery64Details.GetStringValue("InstallDate")
		fmt.Println(installDate)
		if v == 0 {
			continue
		}
		softDetails := Software{displayIcon, displayName, displayVersion, installLocation, publisher, uninstallString, installDate}
		softwareDetails = append(softwareDetails, softDetails)
	}
	return softwareDetails
}
