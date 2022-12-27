package api

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.Engine) {
	r.GET("/saveUser", SaveUser)
	r.GET("/getUserById", GetUserById)
	r.GET("/getAll", GetAll)
	r.GET("/updateById", UpdateById)
	r.GET("/deleteById", DeleteById)
}
