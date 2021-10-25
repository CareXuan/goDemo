package http

import (
	"github.com/gin-gonic/gin"
	"mouse/http/controller"
)

func Route(c *gin.Engine) {
	v1 := c.Group("v1")
	{
		user := v1.Group("user")
		{
			user.POST("/login", controller.LoginIn)
			user.POST("/update", controller.UserUpdate)
			user.GET("/follow", controller.FollowList)
		}
		good := v1.Group("good")
		{
			good.GET("/", controller.GoodList)
			good.POST("/", controller.GoodAdd)
			good.GET("/:id", controller.GetOneGood)
		}
	}
}
