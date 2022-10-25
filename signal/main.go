package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// 监听信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
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
