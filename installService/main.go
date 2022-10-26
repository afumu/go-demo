package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kardianos/service"
)

// 服务自动启动
type program struct {
	log service.Logger
	cfg *service.Config
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	err := s.Stop()
	if err != nil {
		fmt.Println("关闭错误：", err)
	}
	return err
}

func (p *program) run() {
	for {
		now := time.Now()
		sprintf := fmt.Sprintf("%v", now)
		fmt.Println(sprintf)
		time.Sleep(3 * time.Second)

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Println("dir", dir)

		file, err := os.OpenFile(dir+string(filepath.Separator)+"test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0677)
		if err != nil {
			fmt.Println(err)
		}
		file.Write([]byte(sprintf + "\n"))
	}
	wg.Done()
}

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	svcConfig := &service.Config{
		Name:        "agent",
		DisplayName: "agent",
		Description: "态势感知客户端",
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
