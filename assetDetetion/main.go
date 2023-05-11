package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	subnet := "192.168.3."
	startIP := 31
	endIP := 32
	startPort := 3306
	endPort := 3306

	for i := startIP; i <= endIP; i++ {
		// 定义主机的 IP 地址
		host := subnet + strconv.Itoa(i)

		// 启动一个并发任务
		wg.Add(1)
		go func(host string) {
			defer wg.Done()

			// 尝试与指定主机的各个端口建立连接
			for port := startPort; port <= endPort; port++ {
				address := host + ":" + strconv.Itoa(port)
				conn, err := net.Dial("tcp", address)
				if err == nil {
					// 如果连接成功，则说明该端口已开放
					fmt.Printf("Host %s, Port %d: open.\n", host, port)
					conn.Close()
				}
			}
		}(host)
	}

	// 等待所有任务完成
	wg.Wait()
}
