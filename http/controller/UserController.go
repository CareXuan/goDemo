package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mouse/base"
	"mouse/http/model"
	"strconv"
	"strings"
)

var UserGlobalObj model.User

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

func UserUpdate(c *gin.Context) {
	type updateJson struct {
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}
	var updateObj updateJson
	_ = c.BindJSON(&updateObj)
	token := c.GetHeader("token")
	if token == "" {
		base.Forbidden(c, "缺少token", []string{})
		return
	}
	if updateObj.Nickname != "" {
		sql := "UPDATE user SET nickname = ? WHERE token = ?"
		_ = base.Conf.Mysql.QueryRow(sql, updateObj.Nickname, token)
	}
	if updateObj.Password != "" {
		sql := "UPDATE user SET password = ? WHERE token = ?"
		_ = base.Conf.Mysql.QueryRow(sql, updateObj.Password, token)
	}
	base.PostOk(c, "更新成功", []string{})
	return
}

func FollowList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}
	sql := "SELECT user_id,follow_id FROM user_follow WHERE user_id = ? limit ? offset ?"
	res, _ := base.Conf.Mysql.Query(sql, UserGlobalObj.Id, limit, offset)
	var followIds []string
	for res.Next() {
		var userId string
		var followId string
		res.Scan(&userId, &followId)
		followIds = append(followIds, followId)
	}
	if len(followIds) == 0 {
		base.NotFound(c, "您还没有关注其他用户", []model.User{})
		return
	}
	users := make(map[int64]model.User)
	sql = fmt.Sprintf("SELECT id,user_id,nickname,mobile FROM user WHERE id in (%s) and delete_at = 0", strings.Join(followIds, ","))
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var id int64
		var userId int64
		var nickname string
		var mobile int
		res.Scan(&id, &userId, &nickname, &mobile)
		var userObj model.User
		userObj.UserId = userId
		userObj.Nickname = nickname
		userObj.Mobile = mobile
		users[id] = userObj
	}
	sql = fmt.Sprintf("SELECT resource,related_id FROM resource WHERE related_id in (%s) and related_type = 100 and delete_at = 0 and status = 100", strings.Join(followIds, ","))
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var resource string
		var relatedId int64
		res.Scan(&resource, &relatedId)
		userBak := users[relatedId]
		userBak.Avatar = resource
		users[relatedId] = userBak
	}

	var result []model.User
	for _, v := range users {
		result = append(result, v)
	}
	base.GetOk(c, "获取成功", result)
	return
}
