package controller

import (
	"github.com/gin-gonic/gin"
	"mouse/base"
)

func LoginIn(c *gin.Context) {
	loginType := c.Query("type")
	if loginType == "" {
		base.NotFound(c, "缺少必要参数:type", []string{})
	}
	switch loginType {
	case "mobile":
		mobile, _ := c.GetPostForm("mobile")
		code, _ := c.GetPostForm("code")
		break
	case "password":
		mobile, _ := c.GetPostForm("mobile")
		password, _ := c.GetPostForm("password")
		break
	}

	//fmt.Print(c.Params.ByName("aaa"))
}
