package interview2

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test1(t *testing.T) {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := slice[2:5] // 2, 3, 4
	s2 := s1[2:6:7]  // 4 5 6 7

	s2 = append(s2, 100) // 2, 3, 4, 5, 100
	s2 = append(s2, 200) // 2, 3, 4, 5, 100,200

	s1[2] = 20

	fmt.Println(s1)    // 2, 3, 20
	fmt.Println(s2)    // 2, 3, 20, 5, 100,200
	fmt.Println(slice) // 2, 3, 20, 5, 100,200
}

func Test2(t *testing.T) {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 9)
	x := append(s, 11)
	y := append(s, 12)
	fmt.Println(s, x, y)
}

// 演示不同方式添加，增长不一样
func Test3(t *testing.T) {
	s := []int{1, 2}
	s = append(s, 4, 5, 6)
	//s = append(s, 5)
	//s = append(s, 6)
	fmt.Printf("len=%d, cap=%d", len(s), cap(s))
}

// --------------------------------------------------------------------------------
// 演示内存对齐
func Test4(t *testing.T) {
	s := []string{"1", "2", "3"}
	s2 := []string{"1"}
	fmt.Println(unsafe.Sizeof(s))  // 24
	fmt.Println(unsafe.Sizeof(s2)) // 24
}

type Example struct {
	a bool   // 1个字节
	b int    // 8个字节
	c string // 16个字节
}

type A struct {
	a int32
	b int64
	c int32
}

type B struct {
	a int32
	b int32
	c int64
}

func Test5(t *testing.T) {
	//fmt.Println(unsafe.Sizeof(Example{})) // 32
	//fmt.Println(unsafe.Sizeof(A{}))       // 24
	//fmt.Println(unsafe.Sizeof(B{}))       // 16
	//fmt.Println(unsafe.Alignof(int(1)))
	fmt.Println(unsafe.Sizeof(int64(1)))

}

func Test6(t *testing.T) {
	var s []int
	s = append(s, 1)
	s = append(s, 1)
	s = append(s, 1)
	fmt.Println(len(s), cap(s))
	appendData(s)
	fmt.Println(s)
}

func appendData(s []int) {
	// 这里其实还是对同一个底层数组进行了添加，因为容量没有改变
	// 如果容量发生了改变，那么和外部的切片对应的底层数组就不是同一个数组了
	s = append(s, 10)
	fmt.Println(s)
	fmt.Println(len(s), cap(s))
}
