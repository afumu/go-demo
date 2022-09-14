package osinfo

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/shirou/gopsutil/winservices"
)

func GetCpuInfo() {

	//p, _ := cpu.Percent(time.Second, false)
	//fmt.Printf("CPU:%.1f%%", p[0])

	// data := [][]string{}
	// p_Percent, _ := cpu.Percent(time.Second, false)
	// p_count, _ := cpu.Counts(true)

	// table := tabwriter.NewWriter(os.Stdout)
	// table.SetHeader([]string{"项目", "数量", "百分比"})

	// data = append(data, []string{"CPU", fmt.Sprintf("%d核", p_count),
	//     fmt.Sprintf("%.1f%%", p_Percent[0])})
	// for _, v := range data {
	//     table.Append(v)
	// }
	// table.Render()

	info, _ := load.Avg()
	//data = append(data, []string{"CPU", fmt.Sprintf("%d核", v.Total/1024/1024/1024+1)})
	// CPU负载信息
	fmt.Printf("CPU当前负载信息:%v\n", info)

	cpuInfo, _ := cpu.Info() //总体信息
	fmt.Println(cpuInfo)
}

func GetHostInfo() {
	nInfo, _ := host.Info()
	fmt.Printf("HostID:%v\n", nInfo.HostID)
	fmt.Printf("Hostname:%v\n", nInfo.Hostname)
	fmt.Printf("OS:%v\n", nInfo.OS)
	fmt.Printf("Platform:%v\n", nInfo.Platform)
	fmt.Printf("PlatformFamily:%v\n", nInfo.PlatformFamily)
	fmt.Printf("PlatformVersion:%v\n", nInfo.PlatformVersion)
}

func GetCpuAndMem() {
	data := [][]string{}
	v, _ := mem.VirtualMemory()

	fmt.Printf("全部内存:%dG\n", v.Total/1024/1024/1024)
	fmt.Printf("使用内存:%dG\n", v.Used/1024/1024/1024)
	fmt.Printf("闲置可用内存:%vG\n", v.Available/1024/1024/1024)

	fmt.Printf("已经使用内存:%.1f%%\n", v.UsedPercent)
	fmt.Printf("已使用百分比:%.2f", v.UsedPercent)

	//需要转型
	// num, _ := strconv.ParseFloat(fmt.Sprintf("已使用百分比:%.2f", v.UsedPercent), 64)
	// fmt.Printf(num)

	data = append(data, []string{"MEM-Available", fmt.Sprintf("%dG", v.Available/1024/1024/1024)})
	data = append(data, []string{"MEM-Total", fmt.Sprintf("%dG", v.Total/1024/1024/1024+1)})
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"项目", "数量"})

	for _, v := range data {
		table.Append(v)

	}
	table.Render()

	fmt.Printf("\r已经使用内存:%.1f%%\n", v.UsedPercent)
	fmt.Printf("\r闲置可用内存:%vG\n", v.Available/1024/1024/1024)
}

func GetDiskInfo() {

	parts, _ := disk.Partitions(true)
	for _, part := range parts {
		fmt.Printf("part:%v\n", part.String())
		diskInfo, _ := disk.Usage(part.Mountpoint)
		fmt.Printf("已使用:%v\n 可使用:%v\n 磁盘总量:%v\n", diskInfo.Used/1024/1024/1024, diskInfo.Free/1024/1024/1024, diskInfo.Total/1024/1024/1024)
	}

}

func GetNetInfo() {
	info, _ := net.IOCounters(true)
	for index, v := range info {
		fmt.Printf("网卡使用情况: %v:%v\n send:%v\n recv:%v\n", index, v, v.BytesSent, v.BytesRecv)

	}
}

func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr
		if len(macAddr) == 0 {
			continue
		}
		fmt.Printf("mac地址：:%v\n", macAddr)
		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

func GetIPAddrs() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
	}

	for _, netInterface := range netInterfaces {
		ipAddr := netInterface.Addrs
		if len(ipAddr) == 0 {
			continue
		}
		fmt.Printf("ip地址：:%v\n", ipAddr)
	}

}

func NetConnectInfo() {
	infoList, _ := net.Connections("all") //可填入tcp、udp、tcp4、udp4等等
	for index, info := range infoList {
		fmt.Printf("-------%v-------\n", index)
		fmt.Println(info)
	}

}

func GetProcess() {
	processList, err := process.Processes()
	if err != nil {
		fmt.Printf("fail to get net processList: %v", err)
	}
	for index, p := range processList {
		fmt.Printf("--------begin %v--------\n", index)
		fmt.Printf("Pid:%v\n", p.Pid)
		pName, _ := p.Name()
		fmt.Printf("Name:%s\n", pName)
		pCmd, _ := p.Cmdline()
		fmt.Printf("Cmdline:%s\n", pCmd)
		pCreateTime, _ := p.CreateTime()
		fmt.Printf("CreateTime:%d\n", pCreateTime)
		fmt.Printf("--------end--------\n")
	}
}

func GetWinService() {
	serviceList, err := winservices.ListServices()
	if err != nil {
		fmt.Printf("fail to get net serviceList: %v", err)
	}
	for index, s := range serviceList {
		fmt.Printf("--------begin %v--------\n", index)
		fmt.Printf("Name:%v\n", s.Name)
		fmt.Printf("DisplayName:%v\n", s.Config.DisplayName)
		fmt.Printf("Status:%v\n", s.Status)
		fmt.Printf("BinaryPathName:%v\n", s.Config.BinaryPathName)
		fmt.Printf("--------end--------\n")
	}
}
