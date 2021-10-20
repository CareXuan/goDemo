package controller

import (
	"github.com/gin-gonic/gin"
	"mouse/base"
	"mouse/http/model"
	"strconv"
	"strings"
)

func GoodList(c *gin.Context) {
	sql := "select id,name,price,intro from good where delete_at = 0 limit 10"
	res, _ := base.Conf.Mysql.Query(sql)
	goods := make(map[int64]model.Good)
	result := []model.Good{}
	goodIds := []string{}
	for res.Next() {
		var id int64
		var name string
		var price float64
		var intro string
		var good model.Good
		res.Scan(&id, &name, &price, &intro)
		good.Name = name
		good.Price = price
		good.Intro = intro
		goods[id] = good
		goodIds = append(goodIds, strconv.FormatInt(id, 10))
	}
	sql = "select resource,related_id from resource where related_type = 200 and related_id in (?) and delete_at = 0"
	res, _ = base.Conf.Mysql.Query(sql, strings.Join(goodIds, ","))
	for res.Next() {
		var resource string
		var relatedId int64
		res.Scan(&resource, &relatedId)
		newGood := goods[relatedId]
		newGood.Resource = append(newGood.Resource, resource)
		result = append(result, newGood)
	}

	base.GetOk(c, "获取成功", result)
}
