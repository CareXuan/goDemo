package controller

import (
	"github.com/gin-gonic/gin"
	"mouse/base"
	"mouse/http/model"
)

func LoginIn(c *gin.Context) {
	loginType := c.Query("type")
	var userObj model.User
	if loginType == "" {
		base.NotFound(c, "缺少必要参数:type", []string{})
	}
	switch loginType {
	case "mobile":
		mobile, _ := c.GetPostForm("mobile")
		code, _ := c.GetPostForm("code")
		loginRes, resultMsg := checkSms(mobile, code)
		if resultMsg != "" {
			base.Forbidden(c, resultMsg, []string{})
			return
		}
		userObj = loginRes
		break
	case "password":
		mobile, _ := c.GetPostForm("mobile")
		password, _ := c.GetPostForm("password")
		loginRes, resultMsg := checkPwd(mobile, password)
		if resultMsg != "" {
			base.Forbidden(c, resultMsg, []string{})
			return
		}
		userObj = loginRes
		break
	default:
		base.Forbidden(c, "不合法的参数:type", []string{})
		return
	}

	if userObj.UserId != 0 {
		base.PostOk(c, "登录成功", userObj)
	}

}

func checkSms(mobile string, code string) (model.User, string) {
	var loginCount int
	var userObj model.User
	sql := "select count(*) as count from sms where mobile = ? and code = ? and delete_at = 0"
	_ = base.Conf.Mysql.QueryRow(sql, mobile, code).Scan(&loginCount)
	if loginCount <= 0 {
		return userObj, "请先获取验证码或验证码输入错误"
	}
	var userCount int
	sql = "select count(*) as count from user where mobile = ? and delete_at = 0"
	_ = base.Conf.Mysql.QueryRow(sql, mobile).Scan(&userCount)
	if userCount > 0 {
		userObj = model.GetOneUserByMobile(mobile)
	} else {
		userObj = model.CreateOneUser(mobile)
	}
	return userObj, ""
}

func checkPwd(mobile string, password string) (model.User, string) {
	var userCount int
	var userObj model.User
	sql := "select count(*) as count from user where mobile = ? and delete_at = 0"
	_ = base.Conf.Mysql.QueryRow(sql, mobile).Scan(&userCount)
	if userCount < 1 {
		return userObj, "您尚未完成注册，无法进行密码登录"
	}
	var loginCount int
	sql = "select count(*) as count from user where mobile = ? and password = ? and delete_at = 0"
	_ = base.Conf.Mysql.QueryRow(sql, mobile, password).Scan(&loginCount)
	if loginCount < 1 {
		return userObj, "密码错误，请重新输入"
	} else {
		userObj = model.GetOneUserByMobile(mobile)
	}
	return userObj, ""
}
