package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	time.Sleep(2 * time.Second)
	end := time.Now()
	start = start.Add(10 * time.Second)

	fmt.Println(start.After(end))

}
