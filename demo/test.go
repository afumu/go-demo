package main

import "fmt"

type User struct {
	Name string
}

func setUser(list []*User) {
	for _, v := range list {
		v.Name = "lisi"
	}
}

func main() {
	const AgentServiceName = "hostAgent"
	const AgentServiceDisplayName = "Agent"
	const AgentServiceDesc = "态势感知主机代理"

	const UpgradeServiceName = "agentUpgrade"
	const UpgradeServiceDisplayName = "Agent Upgrade"
	const UpgradeServiceDesc = "态势感知主机代理升级服务"

	// 检测agent状态
	const CheckAgentStatus = `
#!/bin/bash
systemctl status ` + AgentServiceName + `
exit 0
`
	fmt.Println(CheckAgentStatus)
}
