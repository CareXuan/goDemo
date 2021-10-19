package base

import "github.com/gin-gonic/gin"

func GetOk(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  msg,
		"data": data,
	})
}

func PostOk(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"code": 1000,
		"msg":  msg,
		"data": data,
	})
}

func NotFound(c *gin.Context, msg string, data interface{}) {
	c.JSON(404, gin.H{
		"code": 2000,
		"msg":  msg,
		"data": data,
	})
}

func Forbidden(c *gin.Context, msg string, data interface{}) {
	c.JSON(403, gin.H{
		"code": 2001,
		"msg":  msg,
		"data": data,
	})
}
