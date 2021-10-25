package model

import (
	"mouse/base"
)

type Good struct {
	Id       int64    `json:"-"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Intro    string   `json:"intro"`
	Type     int      `json:"-"`
	UserId   int64    `json:"-"`
	User     User     `json:"user"`
	Status   int      `json:"-"`
	Resource []string `json:"resource"`
	CreateAt int      `json:"-"`
	UpdateAt int      `json:"-"`
	DeleteAt int      `json:"-"`
}

func GetOneContentGood(goodId string) (Good, string) {
	sql := "SELECT id,name,price,intro,user_id FROM good WHERE delete_at = 0 AND id = ?"
	res := base.Conf.Mysql.QueryRow(sql, goodId)
	var goodObj Good
	res.Scan(&goodObj.Id, &goodObj.Name, &goodObj.Price, &goodObj.Intro, &goodObj.UserId)
	if goodObj.Id == 0 {
		return goodObj, "当前商品不存在"
	}
	sql = "SELECT resource FROM resource WHERE related_type = 200 AND related_id = ? AND status = 100 AND delete_at = 0"
	resourceRes, _ := base.Conf.Mysql.Query(sql, goodId)
	for resourceRes.Next() {
		var resource string
		resourceRes.Scan(&resource)
		goodObj.Resource = append(goodObj.Resource, resource)
	}
	sql = "SELECT nickname,mobile FROM user WHERE id = ?"
	_ = base.Conf.Mysql.QueryRow(sql, goodObj.UserId).Scan(&goodObj.User.Nickname, &goodObj.User.Mobile)

	sql = "SELECT resource FROM resource WHERE related_id = ? AND related_type = 100 AND delete_at = 0 AND status = 100"
	_ = base.Conf.Mysql.QueryRow(sql, goodObj.UserId).Scan(&goodObj.User.Avatar)

	return goodObj, ""
}
