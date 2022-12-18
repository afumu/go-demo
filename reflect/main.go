package main

import (
	"fmt"
	"reflect"
)

type Address struct {
	Name string `json:"name"`
	Size int
}

type User struct {
	Username string `json:"username" id:"123"`
	Age      int
	Address  *Address `json:"address1"`
}

type Iter interface {
}

func main() {

	// 1.在已经知道类型的情况，可以使用
	/*	var num float64 = 3.1415926
		pointer := reflect.ValueOf(&num)
		value := reflect.ValueOf(num)
		a := pointer.Interface().(*float64)
		fmt.Println(a)
		f := value.Interface().(float64)
		fmt.Println(f)
		time.Sleep(3 * time.Second)
		fmt.Println("--------------")*/

	//var dest *User = &User{}
	var i Iter
	//var i2 Iter
	//i2 = dest
	i = 1
	of := reflect.TypeOf(i)
	fmt.Println(of.Kind())
	//i = dest
	//value := reflect.ValueOf(dest)
	//fmt.Println(value.Kind() == reflect.Interface)
	//if value.Kind() == reflect.Ptr && value.IsNil() {
	//	value = reflect.New(value.Type().Elem())
	//}

	//print(i)

	//fmt.Println(i)

	// 2.未知类型的情况

	//var user *User
	//user = &User{Username: "aa"}
	//valueOf := reflect.ValueOf(user)
	//i := valueOf.Interface()
	//fmt.Println())

	//var iter Iter
	//iter = 1
	//print(iter)

	// TypeOf 获取对象的具体类型
	/*	typeOf := reflect.TypeOf(user)
		kind := typeOf.Kind()
		fmt.Println(kind.String())

		zero := reflect.Zero(typeOf.Elem())
		fmt.Println(zero)*/

	// 如果是一个指针类型的话，需要通过elem方法来获取类型名称
	// 注意数组，数组，map，slice，interface都是单独的类型
	/*	if kind == reflect.Pointer {
			fmt.Println(typeOf.Elem().Name())
		} else {
			fmt.Println(typeOf.Name())
		}
	*/
	// 遍历结构体
	/*	fields := typeOf.Elem().NumField()
		fmt.Println("---------------------------")
		for i := 0; i < fields; i++ {
			field := typeOf.Elem().Field(i)
			fmt.Println(field.Tag.Get("id"))
		}*/

	//if name, ok := typeOf.Elem().FieldByName("Address"); ok {
	//	get := name.Tag.Get("json")
	//	fmt.Println(get)
	//}

}

func print(dest interface{}) {
	//value := reflect.ValueOf(dest)
	of := reflect.TypeOf(dest)
	//modelType := value.Type()
	//modelType := reflect.Indirect(value).Type()

	// 如果 modelType 是一个接口的话
	fmt.Println(of.Kind() == reflect.Interface)
	//if value.Kind() == reflect.Interface {
	//	//value = reflect.Indirect(reflect.ValueOf(dest)).Elem().Type()
	//}
	//of := reflect.ValueOf()
	//fmt.Println(of.Elem())
}
