package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"sync"
)

var query64Path = "Software\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
var queryPath = "Software\\Microsoft\\Windows\\CurrentVersion\\Uninstall"

func QueryWindowsSoftware() {
	queryKey := func(w *sync.WaitGroup) {
		defer w.Done()

		// 打开注册表文件
		kQuery64, err := registry.OpenKey(registry.LOCAL_MACHINE, query64Path, registry.READ)
		if err != nil {
			//logUtil.Error(err)
			return
		}

		// 读取注册表所有的key
		keyNames, err := kQuery64.ReadSubKeyNames(0)
		if err != nil {
			//logUtil.Error(err)
			return
		}

		// 查询出query64Path下面的程序详情，并且添加到SoftwareDetails
		getSoftwareDetails(keyNames, query64Path)
		//*res = append(*res, softwareDetailsSli...)

		k, err1 := registry.OpenKey(registry.LOCAL_MACHINE, queryPath, registry.READ)
		if err1 != nil {
			//logUtil.Error(err1)
			return
		}

		// 读取注册表所有的key
		keyNames, err1 = k.ReadSubKeyNames(0)
		if err1 != nil {
			//logUtil.Error(err1)
			return
		}

		//*res = append(*res, getSoftwareDetails(keyNames, queryPath)...)
	}

	//var res []types.WinSoftware
	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(1)

	go queryKey(waitGroup)
	waitGroup.Wait()
	//return nil
}

func getSoftwareDetails(appName []string, queryPath string) {
	//softwareDetails := []types.WinSoftware{}
	for _, value := range appName {
		kQuery64Details, err := registry.OpenKey(registry.LOCAL_MACHINE, queryPath+"\\"+value, registry.READ)
		if err != nil {
			continue
		}
		marshal, _ := json.Marshal(kQuery64Details)
		fmt.Println(string(marshal))
		displayIcon, _, err := kQuery64Details.GetStringValue("DisplayIcon")
		Size, _, err := kQuery64Details.GetStringValue("DisplaySize")
		fmt.Println(displayIcon)
		fmt.Println(Size)
		//displayName, v, err := kQuery64Details.GetStringValue("DisplayName")
		//displayVersion, _, err := kQuery64Details.GetStringValue("DisplayVersion")
		installLocation, _, err := kQuery64Details.GetStringValue("InstallLocation")
		fmt.Println(installLocation)
		//publisher, _, err := kQuery64Details.GetStringValue("Publisher")
		//uninstallString, _, err := kQuery64Details.GetStringValue("UninstallString")
		//installDate, _, err := kQuery64Details.GetStringValue("InstallDate")
		//if v == 0 {
		//	continue
		//}
		//softDetails := types.WinSoftware{displayIcon, displayName, displayVersion, installLocation, publisher, uninstallString, installDate}
		//softwareDetails = append(softwareDetails, softDetails)
	}
	//return softwareDetails
}
