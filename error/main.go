package main

import (
	"errors"
	"fmt"
	"runtime"
)

func WrapError(wrapMsg string, err error) error {
	pc, file, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return errors.New("WrapError 方法获取堆栈失败")
	}
	if err == nil {
		return nil
	} else {
		errMsg := fmt.Sprintf("%s \n\tat %s:%d (Method %s)\nCause by: %s\n", wrapMsg, file, line, f.Name(), err.Error())
		return errors.New(errMsg)
	}
}

func GetErrorStack(preStr string, err error) string {
	pc, file, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	if !ok {
		return "GetErrorStack 方法获取堆栈失败，返回错误原信息：" + err.Error()
	}
	if err == nil {
		return ""
	} else {
		errMsg := fmt.Sprintf("%s \n\tat %s:%d (Method %s)\nCause by: %s\n", preStr, file, line, f.Name(), err.Error())
		return errMsg
	}
}

func TestWrapError() {
	err := errors.New("runtime error: invalid memory address or nil pointer dereference")
	err = WrapError("seek出现异常", err)
	//err = WrapError("twoWrap", err)
	fmt.Print(err)
}

func main() {
	TestWrapError()
}
