package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	//ports := GetOpenPorts()
	//fmt.Println(ports)
	test()
}

func test() {
	str := `
Rule Name:                            Microsoft Edge (mDNS-In-1)
----------------------------------------------------------------------
Enabled:                              Yes
Direction:                            In
Profiles:                             Domain,Private,Public
Grouping:                             Microsoft Edge WebView2 Runtime
LocalIP:                              Any
RemoteIP:                             Any
Protocol:                             UDP
LocalPort:                            5353
RemotePort:                           Any
Edge traversal:                       No
Action:                               Allow

Rule Name:                            Microsoft Edge (mDNS-In-2)
----------------------------------------------------------------------
Enabled:                              Yes
Direction:                            In
Profiles:                             Domain,Private,Public
Grouping:                             Microsoft Edge
LocalIP:                              Any
RemoteIP:                             Any
Protocol:                             UDP
LocalPort:                            5353
RemotePort:                           Any
Edge traversal:                       No
Action:                               Allow

Rule Name:                            Microsoft Teams
----------------------------------------------------------------------
`
	re := regexp.MustCompile(`(?s)(Rule Name:.*?)Rule Name:`)
	//re := regexp.MustCompile(`(?s)Rule Name:\s+(.+?)\n-+\n((?s:.+?)\n\n|$)`)
	matches := re.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		fmt.Println(match[1])
	}

}

func GetOpenPorts() string {
	if runtime.GOOS != "windows" {
		return getLinuxOpenPorts()
	}
	return getWindowOpenPorts()
}

func getLinuxOpenPorts() string {
	cmd := exec.Command("firewall-cmd", "--list-ports")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	portsStr := string(out)
	portsStr = strings.ReplaceAll(portsStr, " ", ",")
	portsStr = strings.ReplaceAll(portsStr, "/tcp", "")
	portsStr = strings.ReplaceAll(portsStr, "/udp", "")
	return portsStr
}

func getWindowOpenPorts() string {
	cmd := exec.Command("cmd", "/c", "chcp 437 & netsh advfirewall firewall show rule name=all dir=in ")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return ""
	}
	str := string(stdout.Bytes())
	re := regexp.MustCompile(`LocalPort:\s+(\S+)`)
	matches := re.FindAllStringSubmatch(str, -1)
	// 输出匹配的本地端口值
	var openPorts []string
	var openPortsMap = make(map[string]string)
	for _, match := range matches {
		matchStr := match[1]
		matchStr = strings.Trim(matchStr, "")
		matchStr = strings.ReplaceAll(matchStr, " ", "")
		split := strings.Split(matchStr, ",")
		for _, v := range split {
			_, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			openPortsMap[v] = ""
		}
	}
	for key := range openPortsMap {
		openPorts = append(openPorts, key)
	}
	return strings.Join(openPorts, ",")
}

// 用go读取window操作系统，遍历出所有的已经安装的软件
