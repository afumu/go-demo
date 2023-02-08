package main

import (
	"fmt"
	"time"
)

func main() {

	// chan 实现定时
	//after()

	// 实现定时任务
	tick := time.Tick(3 * time.Second)
	for {
		select {
		case <-tick:
			fmt.Println("111")
		}
	}
}

func after() {
	afterChan := time.After(3 * time.Second)
	select {
	case <-afterChan:
		fmt.Println("111")
	}
}
