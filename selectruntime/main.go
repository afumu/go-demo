package main

import "fmt"

func main() {
	var c = make(chan struct{})
	select {
	case <-c:
		fmt.Println("读取成功")
	default:
		fmt.Println("default")
	}

}