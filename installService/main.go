package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kardianos/service"
)

type program struct {
	log service.Logger
	cfg *service.Config
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	for {
		fmt.Println("1111111111111111")
		time.Sleep(3 * time.Second)
	}
	wg.Done()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	svcConfig := &service.Config{
		Name:        "test-service",
		DisplayName: "test-service",
		Description: "test-service",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			x := s.Install()
			if x != nil {
				fmt.Println("error:", x.Error())
				return
			}
			fmt.Println("服务安装成功")
			return
		} else if os.Args[1] == "uninstall" {
			x := s.Uninstall()
			if x != nil {
				fmt.Println("error:", x.Error())
				return
			}
			fmt.Println("服务卸载成功")
			return
		}
	}
	err = s.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	wg.Wait()
}
