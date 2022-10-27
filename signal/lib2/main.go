package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

func main() {
	file, err := os.ReadFile("pid.txt")
	if err != nil {
		fmt.Println(err)
	}
	pid, err := strconv.Atoi(string(file))
	if err != nil {
		fmt.Println(err)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println(err)
	}

	err = process.Signal(syscall.SIGWINCH)

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)

}
