package base

import "github.com/gin-gonic/gin"

func GetOk(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  msg,
		"data": data,
	})
}
