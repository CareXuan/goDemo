package http

import (
	"github.com/gin-gonic/gin"
	"mouse/http/api"
)

func Route(c *gin.Engine) {
	c.GET("/ttt", api.Test1)
}
