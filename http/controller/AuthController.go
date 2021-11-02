package controller

import (
	"github.com/gin-gonic/gin"
	"mouse/base"
	"mouse/http/model"
)

//
// LoginIn
// @Description: 登录
// @param c
//
func LoginIn(c *gin.Context) {
	loginType := c.Query("type")
	var userObj model.User
	if loginType == "" {
		base.NotFound(c, "缺少必要参数:type", []string{})
	}
	switch loginType {
	case "mobile":
		type MobileJson struct {
			Mobile string `json:"mobile"`
			Code   string `json:"code"`
		}
		var MobileObj MobileJson
		_ = c.BindJSON(&MobileObj)
		loginRes, resultMsg := checkSms(MobileObj.Mobile, MobileObj.Code)
		if resultMsg != "" {
			base.Forbidden(c, resultMsg, []string{})
			return
		}
		userObj = loginRes
		break
	case "password":
		type PasswordJson struct {
			Mobile   string `json:"mobile"`
			Password string `json:"password"`
		}
		var PasswordObj PasswordJson
		_ = c.BindJSON(&PasswordObj)
		loginRes, resultMsg := checkPwd(PasswordObj.Mobile, PasswordObj.Password)
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

	if userObj.Token != "" {
		c.Header("token", userObj.Token)
		UserGlobalObj = userObj
		base.PostOk(c, "登录成功", userObj)
		return
	} else {
		base.Forbidden(c, "登录失败，请重试", []string{})
		return
	}
}

//
//  checkSms
//  @Description: 检测验证码是否正确
//  @param mobile
//  @param code
//  @return model.User
//  @return string
//
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

//
//  checkPwd
//  @Description: 检查密码登录
//  @param mobile
//  @param password
//  @return model.User
//  @return string
//
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
