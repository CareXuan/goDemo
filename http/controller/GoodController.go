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
	sql := fmt.Sprintf("select id from good where delete_at = 0 limit %d offset %d", limit, offset)
	res, _ := base.Conf.Mysql.Query(sql)
	var goodIds []string
	for res.Next() {
		var id int64
		res.Scan(&id)
		goodIds = append(goodIds, strconv.FormatInt(id, 10))
	}
	if len(goodIds) == 0 {
		base.NotFound(c, "没有更多商品了", []model.Good{})
		return
	}
	goods := model.GetManyContentGood(goodIds)

	var result []model.Good
	for _, v := range goods {
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
