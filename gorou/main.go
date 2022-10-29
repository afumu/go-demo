package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	runtime.GOMAXPROCS(1)
	for i := 0; i < 500000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				//fmt.Printf("协程: %d 执行了,睡眠5s\n", i)
				time.Sleep(5 * time.Second)
				//var sum int
				//for j := 0; j < 100000; j++ {
				//	sum += j
				//}
			}
		}(i)
	}
	fmt.Println("50万协程启动完毕")
	wg.Wait()
	fmt.Println("over")
}
