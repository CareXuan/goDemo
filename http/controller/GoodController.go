package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mouse/base"
	"mouse/http/model"
	"strconv"
	"strings"
)

//
// GoodList
// @Description: 商品列表
// @param c
//
func GoodList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := 0
	if page > 0 {
		offset = (page - 1) * limit
	}
	sql := fmt.Sprintf("select id,name,price,intro,user_id from good where delete_at = 0 limit %d offset %d", limit, offset)
	res, _ := base.Conf.Mysql.Query(sql)
	goods := make(map[int64]model.Good)
	var goodIds []string
	var userIds []string
	for res.Next() {
		var id int64
		var name string
		var price float64
		var intro string
		var userId int64
		var good model.Good
		res.Scan(&id, &name, &price, &intro, &userId)
		good.Name = name
		good.Price = price
		good.Intro = intro
		good.UserId = userId
		good.Resource = []string{}
		goods[id] = good
		userIds = append(userIds, strconv.FormatInt(userId, 10))
		goodIds = append(goodIds, strconv.FormatInt(id, 10))
	}
	//sql = "select resource,related_id from resource where related_type = 200 and related_id in (?) and status = 100 and delete_at = 0"
	sql = fmt.Sprintf("select resource,related_id from resource where related_type = 200 and related_id in (%s) and status = 100 and delete_at = 0", strings.Join(goodIds, ","))
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var resource string
		var relatedId int64
		res.Scan(&resource, &relatedId)
		newGood := goods[relatedId]
		newGood.Resource = append(newGood.Resource, resource)
		goods[relatedId] = newGood
	}

	sql = fmt.Sprintf("SELECT id,nickname,mobile FROM user WHERE id IN (%s)", strings.Join(userIds, ","))
	res, _ = base.Conf.Mysql.Query(sql)
	users := make(map[int64]model.User)
	for res.Next() {
		var id int64
		var nickname string
		var mobile int
		res.Scan(&id, &nickname, &mobile)
		var user2 model.User
		user2.Nickname = nickname
		user2.Mobile = mobile
		users[id] = user2
	}

	sql = fmt.Sprintf("SELECT resource,related_id FROM resource WHERE related_id in (%s) AND related_type = 100 AND delete_at = 0 AND status = 100", strings.Join(userIds, ","))
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var resource string
		var relatedId int64
		res.Scan(&resource, &relatedId)
		userBak := users[relatedId]
		userBak.Avatar = resource
		users[relatedId] = userBak
	}

	var result []model.Good
	for _, v := range goods {
		goodUser := users[v.UserId]
		v.User = goodUser
		result = append(result, v)
	}

	base.GetOk(c, "获取成功", result)
	return
}

//
// GoodAdd
// @Description: 发布商品
// @param c
//
func GoodAdd(c *gin.Context) {
	type goodJson struct {
		Name     string   `json:"name"`
		Intro    string   `json:"intro"`
		Price    float64  `json:"price"`
		Pics     []string `json:"pics"`
		GoodType int      `json:"type"`
	}
	var goodObj goodJson
	_ = c.BindJSON(&goodObj)
	sql := "INSERT INTO good(name,intro,price,type,user_id) VALUES(?,?,?,?,?)"
	stmt, _ := base.Conf.Mysql.Prepare(sql)
	//TODO:给userid改了
	res, _ := stmt.Exec(goodObj.Name, goodObj.Intro, goodObj.Price, goodObj.GoodType, UserGlobalObj.Id)
	insertId, _ := res.LastInsertId()
	var picValuesArr []string
	for _, v := range goodObj.Pics {
		picValuesArr = append(picValuesArr, fmt.Sprintf("(%d,%d,'%s')", 200, insertId, v))
	}
	sql = "INSERT INTO resource(related_type,related_id,resource) VALUES" + strings.Join(picValuesArr, ",")
	_ = base.Conf.Mysql.QueryRow(sql)
	base.PostOk(c, "发布成功", []string{})
	return
}

//
// GetOneGood
// @Description: 获取一个商品
// @param c
//
func GetOneGood(c *gin.Context) {
	goodId := c.Param("id")
	if goodId == "" {
		base.Forbidden(c, "缺少参数ID", []string{})
		return
	}
	goodObj, msg := model.GetOneContentGood(goodId)
	if msg != "" {
		base.Forbidden(c, msg, goodObj)
	}
	base.GetOk(c, "获取成功", goodObj)
	return
}
