package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 非阻塞
func Test1(t *testing.T) {
	var stopSignal = make(chan struct{})
	// 这种非阻塞的发送信号
	// 如果在select中出现了default，那么编译器任务这是一个非阻塞的信号
	go func() {
		for {
			select {
			case <-stopSignal:
				fmt.Println("收到关闭信号")
				return
			default:
				fmt.Println("睡眠2秒")
				time.Sleep(2 * time.Second)
			}
		}
	}()
	time.Sleep(5 * time.Second)
	stopSignal <- struct{}{}
	fmt.Println("main closed")
}

// 阻塞
func Test2(t *testing.T) {
	var stopSignal = make(chan struct{})
	var stopSignal2 = make(chan struct{})
	go func() {
		for {
			select {
			// 如果没有收到信号，这里会一种阻塞
			case <-stopSignal:
				fmt.Println("收到关闭信号1")
			case <-stopSignal2:
				fmt.Println("收到关闭信号2")
				return
			}
		}
	}()
	time.Sleep(2 * time.Second)
	stopSignal <- struct{}{}

	time.Sleep(2 * time.Second)
	stopSignal2 <- struct{}{}

	time.Sleep(2 * time.Second)
	close(stopSignal)
	close(stopSignal2)
	fmt.Println("关闭channel")
	time.Sleep(2 * time.Second)
	fmt.Println("main closed")
}

// context的使用
// 通过context 来控制协程的上下文
func Test3(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancel")
				return
			default:
				fmt.Println("sleep 2 s")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("send cancel signal")
	cancelFunc()
	time.Sleep(2 * time.Second)
	fmt.Println("main closed")
}

// context的使用
// 通过context 来控制协程的上下文
// 通过context 来通知所有的子协程
func Test4(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		go func() {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("子协程 cancel")
					return
				default:
					fmt.Println("子协程 sleep 3 s")
					time.Sleep(3 * time.Second)
				}
			}
		}()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("父协程 cancel")
				return
			default:
				fmt.Println("父协程 sleep 2 s")
				time.Sleep(2 * time.Second)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("send cancel signal")
	cancelFunc()
	time.Sleep(2 * time.Second)
	fmt.Println("main closed")
}

// 通过context 来通知所有的协程
func Test5(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go watch(ctx, "监控1")
	go watch(ctx, "监控2")
	go watch(ctx, "监控3")
	time.Sleep(4 * time.Second)
	fmt.Println("准备关闭所有的协程")
	cancelFunc()
	time.Sleep(1 * time.Second)
	fmt.Println("main 关闭")
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name + " cancel")
			return
		default:
			fmt.Println(name + " sleep 2 s")
			time.Sleep(2 * time.Second)
		}
	}
}

func Test6(t *testing.T) {
	now := time.Now()
	deadline := now.Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	go watch2(ctx, "监控1")
	go watch2(ctx, "监控2")

	fmt.Println("现在开始等待5秒,time=", time.Now().Unix())
	time.Sleep(5 * time.Second)
	fmt.Println("等待5秒结束,准备调用cancel()函数，发现两个子协程已经结束了，time=", time.Now().Unix())
	cancel()
}

// 单独的监控协程
func watch2(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "收到信号，监控退出,time=", time.Now().Unix())
			return
		default:
			fmt.Println(name, "goroutine监控中,time=", time.Now().Unix())
			time.Sleep(1 * time.Second)
		}
	}
}
