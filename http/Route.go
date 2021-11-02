package http

import (
	"github.com/gin-gonic/gin"
	"mouse/http/controller"
	"mouse/http/middleware"
)

func Route(c *gin.Engine) {
	v1 := c.Group("v1")
	{
		auth := v1.Group("auth")
		{
			auth.POST("/login", controller.LoginIn)
		}
		user := v1.Group("user", middleware.TokenCheck)
		{
			user.PUT("/update", controller.UserUpdate)
			user.GET("/follow", controller.FollowList)
			user.GET("/collect", controller.CollectList)
			user.POST("/follow/:target_uid", controller.FollowOne)
			user.POST("/collect/:target_good_id", controller.CollectOne)
			user.DELETE("/follow/:target_uid", controller.UnfollowOne)
			user.DELETE("/collect/:target_good_id", controller.UncollectOne)
			user.POST("/look/:target_good_id", controller.LookOneGood)
			user.GET("/good", controller.UserGood)
			user.GET("/good/buy", controller.BuyList)
			user.GET("/good/sell", controller.SellList)
		}
		good := v1.Group("good")
		{
			good.GET("/", controller.GoodList)
			good.POST("/", controller.GoodAdd)
			good.GET("/:id", controller.GetOneGood)
		}
		trade := v1.Group("trade")
		{
			trade.POST("/:good_id", controller.TradeCreate)
		}
	}
}
