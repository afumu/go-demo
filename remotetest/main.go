package main

import (
	"fmt"
	"time"
)

func main() {
	var age = 18
	var name = "surpass"
	fmt.Println("name is %s, age is %d", name, age)
	for i := 0; i < 10000; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("time======", time.Now())
	}
}
