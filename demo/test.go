package main

import "fmt"

type User struct {
	Name string
}

func setUser(list []*User) {
	for _, v := range list {
		v.Name = "lisi"
	}
}

func main() {
	var list []*User
	list = append(list, &User{Name: "zhangsan"})
	setUser(list)
	for _, v := range list {
		fmt.Println(v)
	}
}
