package main

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
)

// 不支持window，垃圾
func main() {
	gin.SetMode(gin.ReleaseMode)
	r := setRouter()
	if err := endless.ListenAndServe(":9090", r); err != nil {
		log.Fatalf("listen: %s\n", err)
	}
}
func setRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.String(200, "hello world AAA!")
	})
	return r
}
