package main

import (
	"fmt"
	"os"
)

func main() {
	/*	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("执行：", i)
		if i == 5 {
			panic("错误")
		}
	}*/
	strings := os.Args[0:0]
	fmt.Println(strings)
}
