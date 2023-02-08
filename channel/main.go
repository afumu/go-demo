package main

import (
	"fmt"
	"time"
)

func main() {
	c1, c2 := make(chan int), make(chan int)
	go func() {
		time.Sleep(time.Second * 5)
		c1 <- 5
		close(c1)
	}()

	go func() {
		time.Sleep(time.Second * 7)
		c2 <- 7
		close(c2)
	}()

	var ok1, ok2 bool
	for {
		select {
		case x := <-c1:
			ok1 = true
			fmt.Println(x)
		case x := <-c2:
			ok2 = true
			fmt.Println(x)
		}

		if ok1 && ok2 {
			break
		}
	}
	fmt.Println("program end")
}

func print() {
	// 创建一个长度为4，类型为 chan int 的切片
	chSlice := make([]chan int, 4)

	// 遍历切片创建4个协程
	for i, _ := range chSlice {
		chSlice[i] = make(chan int)
		go func(i int) {
			for {
				// 获取当前channel值并打印
				v := <-chSlice[i]
				fmt.Println(v + 1)
				time.Sleep(time.Second)
				// 把下一个值写入下一个channel，等待下一次消费
				chSlice[(i+1)%4] <- (v + 1) % 4
			}
		}(i)
	}
	//      i:   0 1 2 3
	// (i+1)%4:  1 2 3 0
	//       v:  0 1 2 3
	// (v+1)%4:  1 2 3 0
	// 往第一个塞入0
	chSlice[0] <- 0
	select {}
}
