package main

import (
	"fmt"
	"time"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
type App struct {
	// 将封装了代理层的 main.ServiceImpl1 指针注入到 Service 接口，单例模型，需要在标签中指定被注入结构
	ServiceImpl1 Service `singleton:"main.ServiceImpl1"`

	// 将封装了代理层的 main.ServiceImpl2 指针注入到 Service 接口，单例模型，需要在标签中指定被注入结构
	ServiceImpl2 Service `singleton:"main.ServiceImpl2"`

	// 将封装了代理层的 main.ServiceImpl1 指针注入到它的专属接口 'ServiceImpl1IOCInterface'
	// 注入专属接口的命名规则是 '${结构名}IOCInterface'，注入专属接口无需指定被注入结构，标签值为空即可。
	Service1OwnInterface ServiceImpl1IOCInterface `singleton:""`

	// 将结构体指针注入当前字段
	ServiceStruct *ServiceStruct `singleton:""`
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(a.ServiceImpl1.GetHelloString("laurence"))
		fmt.Println(a.ServiceImpl2.GetHelloString("laurence"))

		fmt.Println(a.Service1OwnInterface.GetHelloString("laurence"))

		fmt.Println(a.ServiceStruct.GetString("laurence"))
	}
}

type Service interface {
	GetHelloString(string) string
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl1 struct {
}

func (s *ServiceImpl1) GetHelloString(name string) string {
	return fmt.Sprintf("This is ServiceImpl1, hello %s", name)
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl2 struct {
}

func (s *ServiceImpl2) GetHelloString(name string) string {
	return fmt.Sprintf("This is ServiceImpl2, hello %s", name)
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceStruct struct {
}

func (s *ServiceStruct) GetString(name string) string {
	return fmt.Sprintf("This is ServiceStruct, hello %s", name)
}

func main() {
	// 加载所有结构
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	// 获取结构
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
