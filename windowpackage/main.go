package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strconv"
	"sync"
)

type SoftwareDetails struct {
	DisplayIcon     string
	DisplayName     string
	DisplayVersion  string
	InstallLocation string
	Publisher       string
	UninstallString string
}

func main() {
	exe := QueryWindowsExe()
	fmt.Println(exe)
}

func QueryWindowsExe() []SoftwareDetails {
	softwareDetailsSli := []SoftwareDetails{}

	queryKey := func(w *sync.WaitGroup, startKey registry.Key, res *[]SoftwareDetails) {
		defer w.Done()
		var queryPath string
		var query64Path string
		if strconv.IntSize == 64 {
			query64Path = "Software\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
			queryPath = "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
			kQuery64, err := registry.OpenKey(startKey, query64Path, registry.READ)
			if err != nil {
				return
			}
			keyNames, err := kQuery64.ReadSubKeyNames(0)
			if err != nil {
				return
			}
			//查询出query64Path下面的程序详情，并且添加到SoftwareDetails
			softwareDetailsSli = getSoftwareDetails(startKey, keyNames, query64Path)
			*res = append(*res, softwareDetailsSli...)
		} else {
			queryPath = "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
		}
		k, err1 := registry.OpenKey(startKey, queryPath, registry.READ)
		if err1 != nil {
			return
		}
		// 读取所有子项
		keyNames, err1 := k.ReadSubKeyNames(0)
		if err1 != nil {
			return
		}
		*res = append(*res, getSoftwareDetails(startKey, keyNames, queryPath)...)
	}

	res := []SoftwareDetails{}
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	go queryKey(waitGroup, registry.LOCAL_MACHINE, &res)
	waitGroup.Wait()
	return res
}
func getSoftwareDetails(startKey registry.Key, appName []string, path string) []SoftwareDetails {
	softwareDetails := []SoftwareDetails{}
	for _, value := range appName {
		kQuery64Details, err := registry.OpenKey(startKey, path+"\\"+value, registry.READ)
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
		softDetails := SoftwareDetails{displayIcon, displayName, displayVersion, installLocation, publisher, uninstallString}
		softwareDetails = append(softwareDetails, softDetails)
	}
	return softwareDetails
}
