package controller

import (
	"github.com/gin-gonic/gin"
	"mouse/base"
	"mouse/http/model"
)

//
// TradeCreate
// @Description: 创建订单
// @param c
//
func TradeCreate(c *gin.Context) {
	type TradePrice struct {
		FinishPrice float64 `json:"finish_price"`
	}
	var tradePriceObj TradePrice
	goodId := c.Param("good_id")
	_ = c.BindJSON(&tradePriceObj)
	good, msg := model.GetOneContentGood(goodId)
	if msg != "" {
		base.NotFound(c, msg, []string{})
		return
	}
	if tradePriceObj.FinishPrice == 0 {
		base.NotFound(c, "缺少报价", []string{})
		return
	}
	var tradeCount int
	sql := "SELECT count(*) AS count FROM trade WHERE user_id = ? AND good_id = ? AND delete_at = 0"
	_ = base.Conf.Mysql.QueryRow(sql, UserGlobalObj.Id, goodId).Scan(&tradeCount)
	if tradeCount > 0 {
		base.Forbidden(c, "您已经下单了", []string{})
		return
	}
	sql = "INSERT INTO trade(user_id,target_user_id,good_id,price,finish_price,status) VALUES(?,?,?,?,?,?)"
	stmt, _ := base.Conf.Mysql.Prepare(sql)
	stmt.Exec(UserGlobalObj.Id, good.UserId, good.Id, good.Price, tradePriceObj.FinishPrice, 0)
	base.PostOk(c, "下单成功", []string{})
	return
}
