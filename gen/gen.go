package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Username string
}

func main() {

	var user User = User{
		Username: "zhangan",
	}

	valueOf := reflect.ValueOf(user)
	fmt.Println(valueOf)
	//getUsername(valueOf)

	// 通过调用interface 方法，后续可以进行断言使用
	valueOfInterface := valueOf.Interface()
	fmt.Println(valueOf.Type())

	valueOf.Interface()

	// 这个其实就是类型转换了
	u := valueOfInterface.(User)
	fmt.Println(u)

	getUsername(u)
}

func getUsername(u User) string {
	return u.Username
}
