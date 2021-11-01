package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mouse/base"
	"mouse/http/model"
	"strconv"
	"strings"
	"time"
)

var UserGlobalObj model.User

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

//
// UserUpdate
// @Description: 用户更新
// @param c
//
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
	base.PutOk(c, "更新成功", []string{})
	return
}

//
// FollowList
// @Description: 关注列表
// @param c
//
func FollowList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}
	sql := "SELECT user_id,follow_id FROM user_follow WHERE user_id = ? LIMIT ? OFFSET ? ORDER BY active_at DESC"
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

//
// FollowOne
// @Description: 关注一个用户
// @param c
//
func FollowOne(c *gin.Context) {
	targetUserId := c.Param("target_uid")
	targetUser := model.GetOneUserByUserId(targetUserId)
	userId := UserGlobalObj.Id
	sql := "SELECT count(*) AS count FROM user_follow WHERE user_id = ? AND follow_id = ? AND delete_at = 0"
	var count int
	_ = base.Conf.Mysql.QueryRow(sql, userId, targetUser.Id).Scan(&count)
	if count > 0 {
		base.Forbidden(c, "已经关注过了", []string{})
		return
	}
	sql = "INSERT INTO user_follow(user_id,follow_id) VALUES(?,?)"
	stmt, _ := base.Conf.Mysql.Prepare(sql)
	_, _ = stmt.Exec(userId, targetUser.Id)
	base.PostOk(c, "关注成功", []string{})
	return
}

//
// UnfollowOne
// @Description: 取消关注
// @param c
//
func UnfollowOne(c *gin.Context) {
	targetUserId := c.Param("target_uid")
	targetUser := model.GetOneUserByUserId(targetUserId)
	userId := UserGlobalObj.Id
	sql := "SELECT count(*) AS count FROM user_follow WHERE user_id = ? AND follow_id = ? AND delete_at = 0"
	var count int
	_ = base.Conf.Mysql.QueryRow(sql, userId, targetUser.Id).Scan(&count)
	if count <= 0 {
		base.Forbidden(c, "尚未关注这个用户", []string{})
		return
	}
	sql = "UPDATE user_follow SET delete_at = ? WHERE user_id = ? AND follow_id = ?"
	stmt, _ := base.Conf.Mysql.Prepare(sql)
	_, _ = stmt.Exec(time.Now().Unix(), userId, targetUser.Id)
	base.DeleteOk(c, "取消关注成功", []string{})
	return
}

//
// CollectList
// @Description: 收藏列表
// @param c
//
func CollectList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}
	userId := UserGlobalObj.Id
	sql := "SELECT good_id FROM user_collection WHERE delete_at = 0 AND user_id = " + strconv.FormatInt(userId, 10) + " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	var goodIds []string
	res, _ := base.Conf.Mysql.Query(sql)
	for res.Next() {
		var goodId string
		res.Scan(&goodId)
		goodIds = append(goodIds, goodId)
	}
	if len(goodIds) == 0 {
		base.NotFound(c, "您尚未收藏任何商品", []string{})
		return
	}
	goods := model.GetManyContentGood(goodIds)
	base.GetOk(c, "获取成功", goods)
	return
}

//
// CollectOne
// @Description: 收藏商品
// @param c
//
func CollectOne(c *gin.Context) {
	targetGoodId := c.Param("target_good_id")
	sql := "SELECT count(*) as count FROM good WHERE delete_at = 0 AND id = ?"
	var goodCount int
	_ = base.Conf.Mysql.QueryRow(sql, targetGoodId).Scan(&goodCount)
	if goodCount <= 0 {
		base.NotFound(c, "当前商品不存在", []string{})
		return
	}
	userId := UserGlobalObj.Id
	sql = "SELECT count(*) AS count FROM user_collection WHERE user_id = ? AND good_id = ? AND delete_at = 0"
	var userCount int
	_ = base.Conf.Mysql.QueryRow(sql, userId, targetGoodId).Scan(&userCount)
	if userCount > 0 {
		base.Forbidden(c, "已经收藏过了", []string{})
		return
	}
	sql = "INSERT INTO user_collection(user_id,good_id) VALUES(?,?)"
	stmt, _ := base.Conf.Mysql.Prepare(sql)
	_, _ = stmt.Exec(userId, targetGoodId)
	base.PostOk(c, "收藏成功", []string{})
	return
}

//
// UncollectOne
// @Description: 取消收藏
// @param c
//
func UncollectOne(c *gin.Context) {
	targetGoodId := c.Param("target_good_id")
	sql := "SELECT count(*) as count FROM good WHERE delete_at = 0 AND id = ?"
	var goodCount int
	_ = base.Conf.Mysql.QueryRow(sql, targetGoodId).Scan(&goodCount)
	if goodCount <= 0 {
		base.NotFound(c, "当前商品不存在", []string{})
		return
	}
	userId := UserGlobalObj.Id
	sql = "SELECT count(*) AS count FROM user_collection WHERE user_id = ? AND good_id = ? AND delete_at = 0"
	var userCount int
	_ = base.Conf.Mysql.QueryRow(sql, userId, targetGoodId).Scan(&userCount)
	if userCount <= 0 {
		base.Forbidden(c, "尚未收藏这个商品", []string{})
		return
	}
	sql = "UPDATE user_collection SET delete_at = ? WHERE user_id = ? AND good_id = ?"
	stmt, _ := base.Conf.Mysql.Prepare(sql)
	_, _ = stmt.Exec(time.Now().Unix(), userId, targetGoodId)
	base.DeleteOk(c, "取消收藏成功", []string{})
	return
}

//
// LookOneGood
// @Description: 查看一个商品
// @param c
//
func LookOneGood(c *gin.Context) {
	targetGoodId := c.Param("target_good_id")
	sql := "SELECT count(*) AS count FROM good WHERE delete_at = 0 AND id = ?"
	var goodCount int
	_ = base.Conf.Mysql.QueryRow(sql, targetGoodId).Scan(&goodCount)
	if goodCount <= 0 {
		base.NotFound(c, "当前商品不存在", []string{})
		return
	}
	sql = "SELECT count(*) AS count FROM user_look WHERE delete_at = 0 AND user_id = ? AND good_id = ?"
	var lookCount int
	_ = base.Conf.Mysql.QueryRow(sql, UserGlobalObj.Id, targetGoodId).Scan(&lookCount)
	if lookCount > 0 {
		sql = "UPDATE user_look SET look_at = ? WHERE user_id = ? AND good_id = ? AND delete_at = 0"
		stmt, _ := base.Conf.Mysql.Prepare(sql)
		_, _ = stmt.Exec(time.Now().Unix(), UserGlobalObj.Id, targetGoodId)
	} else {
		sql = "INSERT INTO user_look(user_id,good_id,look_at) VALUES(?,?,?)"
		stmt, _ := base.Conf.Mysql.Prepare(sql)
		_, _ = stmt.Exec(UserGlobalObj.Id, targetGoodId, time.Now().Unix())
	}
	base.PostOk(c, "查看成功", []string{})
	return
}

//
// UserGood
// @Description: 用户商品列表
// @param c
//
func UserGood(c *gin.Context) {
	userId := UserGlobalObj.Id
	sql := "SELECT id FROM good WHERE user_id = ? AND delete_at = 0"
	res, _ := base.Conf.Mysql.Query(sql, userId)
	var goodIds []string
	for res.Next() {
		var id int64
		goodIds = append(goodIds, strconv.FormatInt(id, 10))
	}
	if len(goodIds) == 0 {
		base.NotFound(c, "没有更多商品了", []model.Good{})
		return
	}
	goods := model.GetManyContentGood(goodIds)
	base.GetOk(c, "获取成功", goods)
	return
}

//
// BuyList
// @Description: 我买过的商品
// @param c
//
func BuyList(c *gin.Context) {
	type BuyGoods struct {
		GoodId      int64      `json:"-"`
		Good        model.Good `json:"good"`
		Price       float64    `json:"price"`
		FinishPrice float64    `json:"finish_price"`
		Status      string     `json:"status"`
	}
	userId := UserGlobalObj.Id
	sql := "SELECT good_id,price,finish_price,status FROM trade WHERE user_id = ? AND delete_at = 0"
	res, _ := base.Conf.Mysql.Query(sql, userId)
	var result []BuyGoods
	var goodIds []string
	for res.Next() {
		var goodId int64
		var price float64
		var finishPrice float64
		var status int
		var buyGood BuyGoods
		res.Scan(&goodId, &price, &finishPrice, &status)
		buyGood.GoodId = goodId
		buyGood.Price = price
		buyGood.FinishPrice = finishPrice
		buyGood.Status = model.StatusMapping[status]
		goodIds = append(goodIds, strconv.FormatInt(goodId, 10))
		result = append(result, buyGood)
	}
	goods := model.GetManyContentGoodMapping(goodIds)
	for i := range result {
		result[i].Good = goods[result[i].GoodId]
	}
	base.GetOk(c, "获取成功", result)
}

//
// SellList
// @Description: 我卖出的商品
// @param c
//
func SellList(c *gin.Context) {
	type BuyGoods struct {
		GoodId      int64      `json:"-"`
		Good        model.Good `json:"good"`
		Price       float64    `json:"price"`
		FinishPrice float64    `json:"finish_price"`
		Status      string     `json:"status"`
	}
	userId := UserGlobalObj.Id
	sql := "SELECT good_id,price,finish_price,status FROM trade WHERE targer_user_id = ? AND delete_at = 0"
	res, _ := base.Conf.Mysql.Query(sql, userId)
	var result []BuyGoods
	var goodIds []string
	for res.Next() {
		var goodId int64
		var price float64
		var finishPrice float64
		var status int
		var buyGood BuyGoods
		res.Scan(&goodId, &price, &finishPrice, &status)
		buyGood.GoodId = goodId
		buyGood.Price = price
		buyGood.FinishPrice = finishPrice
		buyGood.Status = model.StatusMapping[status]
		goodIds = append(goodIds, strconv.FormatInt(goodId, 10))
		result = append(result, buyGood)
	}
	goods := model.GetManyContentGoodMapping(goodIds)
	for i := range result {
		result[i].Good = goods[result[i].GoodId]
	}
	base.GetOk(c, "获取成功", result)
}
