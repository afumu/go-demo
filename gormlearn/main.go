package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"test-demo/gormlearn/router"
)

func main() {
	engine := gin.Default()
	router.Init(engine)
	err := engine.Run(":8081")
	if err != nil {
		fmt.Println(err)
	}
}
