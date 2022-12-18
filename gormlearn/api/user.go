package api

import (
	"github.com/gin-gonic/gin"
	"test-demo/gormlearn/dao"
	"time"
)

func SaveUser(c *gin.Context) {
	user := &dao.User{Username: "zcf", Password: "123456", CreateTime: time.Now().UnixMilli()}
	dao.SaveUser(user)
	c.JSON(200, user)
}

func GetUserById(c *gin.Context) {
	user := dao.GetUserById()
	c.JSON(200, user)
}

func GetAll(c *gin.Context) {
	users := dao.GetAll()
	c.JSON(200, users)
}

func UpdateById(c *gin.Context) {
	user := dao.UpdateById()
	c.JSON(200, user)
}

func DeleteById(c *gin.Context) {
	dao.DeleteById()
	c.JSON(200, "")
}
