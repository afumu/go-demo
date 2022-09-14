package main

import "fmt"

type Slice[T int | float32 | float64] []T
type NewType[T int | float32] []T

func main() {

	type MySlice[T int | string | float32] []T
	var mySlice MySlice[int] = []int{1, 2, 3}
	fmt.Printf("%T", mySlice)

	// 以上类型相当于是定义
	type MySlice2 []int

	//type UintSlice[T uint | uint8] MySlice[T]
	//type MySlice2 []int
	//var mySlice2 MySlice2
	//mySlice2 = append(mySlice2, 1)
	//fmt.Println(mySlice)

}
