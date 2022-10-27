package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	go func() {
		pid := os.Getpid()
		fmt.Println("pid:", pid)
		time.Sleep(10 * time.Second)
		file, err := os.OpenFile("version.json", os.O_CREATE|os.O_WRONLY, 0677)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("创建文件成功")
		pidStr := fmt.Sprintf("%d", pid)
		_, err = file.WriteString(pidStr)
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Close()
		fmt.Println("写入pid成功")
	}()

	go func() {
		for {
			now := time.Now()
			message := fmt.Sprintf("%v----------new-----------", now)
			fmt.Println(message)
			time.Sleep(2 * time.Second)
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// 监听信号
	signal.Notify(sigs, syscall.SIGKILL)
	go func() {
		// 接收到信号返回
		sig := <-sigs
		fmt.Println("--------------------------over-------------------")
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("等待信号")
	// 等待信号的接收
	<-done
	fmt.Println("退出")

}
