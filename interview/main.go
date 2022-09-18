package main

import "fmt"

// 需求：通过把一个函数赋值给一个接口
type P interface {
	Add(int, int) int
}

// 中间层
type M func(int, int) int

func (m M) Add(a int, b int) int {
	return m(a, b)
}

func Add(a int, b int) int {
	return a + b
}

// go 面试题
func main() {
	//var p P
	//p = Add(1, 2)
	// 直接赋值是不可以的

	var p P
	p = M(Add)
	r := p.Add(1, 3)
	fmt.Println(r)
}
