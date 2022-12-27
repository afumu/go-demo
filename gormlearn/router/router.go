package router

import (
	"github.com/gin-gonic/gin"
	"test-demo/gormlearn/api"
)

func Init(r *gin.Engine) {
	api.RegisterRouter(r)
}
