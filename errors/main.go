package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	err := my1()
	fmt.Println(err)
	err = errors.Cause(err)
	fmt.Println(err)

	//var s string = `dafdsafdsf="test=123"`
}

func my1() error {
	err := my()
	if err != nil {
		return errors.Wrap(err, "read failed1")
	}
	return nil
}

func my() error {
	err := reader()
	if err != nil {
		return errors.Wrap(err, "read failed")
	}
	return nil
}

func reader() error {
	return errors.New("读取数据错误，路径找不到")
}
