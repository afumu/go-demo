package main

import (
	"fmt"
	"github.com/tangpanqing/aorm/utils"
	"reflect"
)

var FieldMap = make(map[uintptr]string)

type Person struct {
	Id   int
	Name string
	Age  int
}

func main() {
	var p Person
	buildFiledNameMap(&p)
	fmt.Println(FieldMap)
	valueOf := reflect.ValueOf(&p.Age)
	fmt.Println(valueOf.Pointer())
}

func buildFiledNameMap(model any) {
	valueOf := reflect.ValueOf(model)
	typeOf := reflect.TypeOf(model)
	for i := 0; i < valueOf.Elem().NumField(); i++ {
		pointer := valueOf.Elem().Field(i).Addr().Pointer()
		field := typeOf.Elem().Field(i)
		key := utils.UnderLine(field.Name)
		FieldMap[pointer] = key
	}
}
