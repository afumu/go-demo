package main

import (
	"fmt"
)

func main() {
	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop()
	//for range ticker.C {
	//	fmt.Println("ticker tick....")
	//}

	var c1 *Config
	fmt.Println(c1)
	fmt.Println(c1.Path())
}

type Config struct {
	path string
}

func (c *Config) Path() string {
	if c == nil {
		return "/usr/home"
	}
	return c.path
}

func Path(c *Config) string {
	if c == nil {
		return "/usr/home"
	}
	return c.path
}
