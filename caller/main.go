package main

import (
	"fmt"
	"runtime"
)

func main() {
	test()
}

func test() {
	caller, file, line, ok := runtime.Caller(1)
	fmt.Println(caller)
	fmt.Println(file)
	fmt.Println(line)
	fmt.Println(ok)

}
