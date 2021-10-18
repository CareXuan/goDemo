package common

import (
	"github.com/gin-gonic/gin"
	"mouse/api"
)

func Route(c *gin.Engine) {
	c.GET("/ttt", api.Test1)
}
