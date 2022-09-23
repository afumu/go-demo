package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {

	// 获取执行命令
	cmd := exec.Command("D:\\workplace-go\\local\\go-demo\\consoleread\\host\\host.exe")

	cmd.Stdin = os.Stdin
	var wg *sync.WaitGroup
	wg.Add(2)
	// 捕获标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("ERROR:", err)
	}
	readOut := bufio.NewReader(stdout)
	go func() {
		defer wg.Done()
		printOutput(readOut)
	}()

	// 捕获标准错误
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("ERROR:", err)
	}
	readErr := bufio.NewReader(stderr)
	go func() {
		defer wg.Done()
		printOutput(readErr)
	}()

	// 执行命令
	err = cmd.Run()
	if err != nil {
		log.Fatal("ERROR:", err)
	}
	wg.Wait()
}

func printOutput(reader *bufio.Reader) {
	var sumOutput string
	outputBytes := make([]byte, 200)
	for {
		//获取屏幕的实时输出(并不是按照回车分割，所以要结合sumOutput)
		n, err := reader.Read(outputBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			sumOutput += err.Error()
		}
		output := string(outputBytes[:n])
		//输出屏幕内容
		fmt.Print(output)
		sumOutput += output
	}
}
