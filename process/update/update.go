package main

import (
	"bufio"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println(event.Name, event.Op)
				pidStr, err := os.ReadFile("version.json")
				if err != nil {
					fmt.Println(err)
					return
				}
				pid, err := strconv.Atoi(string(pidStr))
				process, err := os.FindProcess(pid)
				if err != nil {
					fmt.Println(err)
				}
				if process != nil {
					process.Signal(syscall.SIGKILL)

					time.Sleep(3 * time.Second)
					fmt.Println("111111111111111")
					exists, _ := PathExists("agent-new.exe")
					if exists {
						err = os.Rename("agent.exe", "agent-old.exe")
						if err != nil {
							fmt.Println(err)
						}

						err = os.Rename("agent-new.exe", "agent.exe")
						go func() {
							executeCommand()
						}()
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("version.json")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
}

func executeCommand() {
	currentPath, err2 := os.Getwd()
	if err2 != nil {
		fmt.Println(err2)
	}
	join := filepath.Join(currentPath, "agent.exe")
	cmd := exec.Command(join)
	printStdout(cmd)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("更新成功")

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}
	if err != nil {
		log.Println(err)
	}
	return false, err
}

func printStdout(cmd *exec.Cmd) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	readOut := bufio.NewReader(stdout)
	go func() {
		printOutput(readOut)
	}()
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
		if n == 0 {
			break
		}
	}
}
