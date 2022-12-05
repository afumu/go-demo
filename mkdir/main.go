package main

import "fmt"

type error interface {
	Error() string
}

type MyInt int
type User struct {
	Age MyInt
}

func main() {
	user := User{}
	sprintf := fmt.Sprintf("%v", user.Age)
	fmt.Println(sprintf)
}
