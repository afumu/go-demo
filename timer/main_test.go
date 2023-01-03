package main

import (
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	var ch = make(chan struct{})
	go func() {
		time.Sleep(3 * time.Second)
		ch <- struct{}{}
	}()

	timer := time.NewTimer(2 * time.Second)

	select {
	case <-ch:
		timer.Stop()
		fmt.Println("收到信号")
	case <-timer.C:
		fmt.Println("超时啦。。。。")
	}
	time.Sleep(7 * time.Second)
}
