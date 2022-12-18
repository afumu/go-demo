package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"test-demo/walker/lib"
	"time"
)

func main() {
	file, err := os.OpenFile("test.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0667)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	sprintf := fmt.Sprintf("%v\n", time.Now())
	file.WriteString(sprintf)
}

func OpenFile(dirname string) (*os.File, error) {
	var method string
	if runtime.GOOS == "windows" {
		method = "Open"
	} else {
		method = "OpenFileNoCache"
	}
	var open = &lib.Open{}
	openValue := reflect.ValueOf(open)
	infoFunc := openValue.MethodByName(method)
	file := infoFunc.Call([]reflect.Value{reflect.ValueOf(dirname)})
	var f *os.File
	if file[0].CanInterface() {
		value := file[0]
		f = value.Interface().(*os.File)
		return f, nil
	}
	return nil, errors.New("读取file失败")
}
