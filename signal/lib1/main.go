package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	pid := os.Getpid()
	fmt.Println("pid", pid)
	file, err := os.OpenFile("pid.txt", os.O_CREATE|os.O_WRONLY, 0667)
	if err != nil {
		fmt.Println(err)
	}
	pidStr := fmt.Sprintf("%d", pid)
	file.WriteString(pidStr)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// 监听信号
	signal.Notify(sigs)
	signal.Ignore(syscall.SIGKILL)
	go func() {
		// 接收到信号返回
		sig := <-sigs
		fmt.Println()
		fmt.Println("--------------------------over-------------------")
		fmt.Println(sig)

		done <- true
	}()
	fmt.Println("awaiting signal")
	// 等待信号的接收
	<-done
	fmt.Println("exiting")

}
