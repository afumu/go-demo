package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	getRandomNumber(1)
}

func getRandomNumber(new int) int {
	// 通过时间戳来设置种子，保证每次运行生成的随机数不同
	rand.Seed(time.Now().UnixNano())

	// 生成 1 到 300 之间的随机数
	randomNumber := rand.Intn(300) + 1

	i := new(struct{})
	fmt.Println(i)
	return randomNumber

}
