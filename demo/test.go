package main

import "test-demo/demo/osinfo"

func Add(x, y int) int {
	return x + y
}

func Subtract(x, y int) int {
	return x - y
}

func main() {
	osinfo.GetHostInfo()
}
