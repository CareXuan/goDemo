package model

import (
	"mouse/base"
	"strconv"
	"strings"
)

type Good struct {
	Id       int64          `json:"-"`
	Name     string         `json:"name"`
	Price    float64        `json:"price"`
	Intro    string         `json:"intro"`
	Type     []RelationType `json:"type"`
	UserId   int64          `json:"-"`
	User     User           `json:"user"`
	Status   int            `json:"-"`
	Resource []string       `json:"resource"`
	CreateAt int            `json:"-"`
	UpdateAt int            `json:"-"`
	DeleteAt int            `json:"-"`
}

//
// GetOneContentGood
// @Description: 获取一个完整的商品
// @param goodId
// @return Good
// @return string
//
func GetOneContentGood(goodId string) (Good, string) {
	sql := "SELECT id,name,price,intro,user_id FROM good WHERE delete_at = 0 AND id = ?"
	var goodObj Good
	_ = base.Conf.Mysql.QueryRow(sql, goodId).Scan(&goodObj.Id, &goodObj.Name, &goodObj.Price, &goodObj.Intro, &goodObj.UserId)
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
	sql = "SELECT type,related_id,type_id FROM relation_type WHERE related_id = ? AND delete_at = 0"
	var typeIds []string
	typeRes, _ := base.Conf.Mysql.Query(sql, goodId)
	for typeRes.Next() {
		var typeNum int
		var relatedId int64
		var typeId int64
		typeRes.Scan(&typeNum, &relatedId, &typeId)
		var typeObj RelationType
		typeObj.Name = RelationMap[typeNum]
		typeObj.RelatedId = relatedId
		typeObj.TypeId = typeId
		typeIds = append(typeIds, strconv.FormatInt(typeId, 10))
		goodObj.Type = append(goodObj.Type, typeObj)
	}

	typeMapping := make(map[int64]string)
	sql = "SELECT id,name FROM type WHERE delete_at = 0 AND id in (" + strings.Join(typeIds, ",") + ")"
	res, _ := base.Conf.Mysql.Query(sql)
	for res.Next() {
		var id int64
		var name string
		res.Scan(&id, &name)
		typeMapping[id] = name
	}
	for i := range goodObj.Type {
		goodObj.Type[i].Value = typeMapping[goodObj.Type[i].TypeId]
	}
	sql = "SELECT nickname,mobile,active_at FROM user WHERE id = ?"
	_ = base.Conf.Mysql.QueryRow(sql, goodObj.UserId).Scan(&goodObj.User.Nickname, &goodObj.User.Mobile, &goodObj.User.ActiveAt)

	sql = "SELECT resource FROM resource WHERE related_id = ? AND related_type = 100 AND delete_at = 0 AND status = 100"
	_ = base.Conf.Mysql.QueryRow(sql, goodObj.UserId).Scan(&goodObj.User.Avatar)

	return goodObj, ""
}

//
// GetManyContentGood
// @Description: 批量获取商品列表
// @param goodIds
// @return []Good
//
func GetManyContentGood(goodIds []string) []Good {
	sql := "SELECT id,name,price,intro,user_id FROM good WHERE delete_at = 0 AND id in (" + strings.Join(goodIds, ",") + ")"
	res, _ := base.Conf.Mysql.Query(sql)
	var userIds []string
	var result []Good
	goods := make(map[int64]Good)
	for res.Next() {
		var id int64
		var name string
		var price float64
		var intro string
		var userId int64
		res.Scan(&id, &name, &price, &intro, &userId)
		var good Good
		good.Id = id
		good.Name = name
		good.Price = price
		good.Intro = intro
		good.UserId = userId
		good.Resource = []string{}
		userIds = append(userIds, strconv.FormatInt(userId, 10))
		goods[id] = good
	}

	sql = "SELECT resource,related_id FROM resource WHERE delete_at = 0 AND related_type = 200 AND status = 100 AND related_id in (" + strings.Join(goodIds, ",") + ")"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var resource string
		var relatedId int64
		res.Scan(&resource, &relatedId)
		newGood := goods[relatedId]
		newGood.Resource = append(newGood.Resource, resource)
		goods[relatedId] = newGood
	}

	sql = "SELECT type,related_id,type_id FROM relation_type WHERE related_id in (" + strings.Join(goodIds, ",") + ") AND delete_at = 0"
	var typeIds []string
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var typeNum int
		var relatedId int64
		var typeId int64
		res.Scan(&typeNum, &relatedId, &typeId)
		var typeObj RelationType
		typeObj.Name = RelationMap[typeNum]
		typeObj.RelatedId = relatedId
		typeObj.TypeId = typeId
		typeIds = append(typeIds, strconv.FormatInt(typeId, 10))
		newGood := goods[relatedId]
		newGood.Type = append(newGood.Type, typeObj)
		goods[relatedId] = newGood
	}

	typeMapping := make(map[int64]string)
	sql = "SELECT id,name FROM type WHERE delete_at = 0 AND id in (" + strings.Join(typeIds, ",") + ")"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var id int64
		var name string
		res.Scan(&id, &name)
		typeMapping[id] = name
	}
	for i := range goods {
		for j := range goods[i].Type {
			goods[i].Type[j].Value = typeMapping[goods[i].Type[j].TypeId]
		}
	}

	users := make(map[int64]User)
	sql = "SELECT id,nickname,mobile,active_at FROM user WHERE id in (" + strings.Join(userIds, ",") + ")"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var id int64
		var nickname string
		var mobile int
		var activeAt int
		res.Scan(&id, &nickname, &mobile, &activeAt)
		var user User
		user.Id = id
		user.Nickname = nickname
		user.Mobile = mobile
		user.ActiveAt = activeAt
		users[id] = user
	}

	sql = "SELECT resource,related_id FROM resource WHERE related_id in (" + strings.Join(userIds, ",") + ") AND related_type = 100 AND delete_at = 0 AND status = 100"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var resource string
		var relatedId int64
		newUser := users[relatedId]
		newUser.Avatar = resource
		users[relatedId] = newUser
	}

	for _, v := range goods {
		v.User = users[v.UserId]
		result = append(result, v)
	}
	return result
}

//
// GetManyContentGoodMapping
// @Description: 批量获取商品列表(map)
// @param goodIds
// @return []Good
//
func GetManyContentGoodMapping(goodIds []string) map[int64]Good {
	sql := "SELECT id,name,price,intro,user_id FROM good WHERE delete_at = 0 AND id in (" + strings.Join(goodIds, ",") + ")"
	res, _ := base.Conf.Mysql.Query(sql)
	var userIds []string
	goods := make(map[int64]Good)
	for res.Next() {
		var id int64
		var name string
		var price float64
		var intro string
		var userId int64
		res.Scan(&id, &name, &price, &intro, &userId)
		var good Good
		good.Id = id
		good.Name = name
		good.Price = price
		good.Intro = intro
		good.UserId = userId
		good.Resource = []string{}
		userIds = append(userIds, strconv.FormatInt(userId, 10))
		goods[id] = good
	}

	sql = "SELECT resource,related_id FROM resource WHERE delete_at = 0 AND related_type = 200 AND status = 100 AND related_id in (" + strings.Join(goodIds, ",") + ")"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var resource string
		var relatedId int64
		res.Scan(&resource, &relatedId)
		newGood := goods[relatedId]
		newGood.Resource = append(newGood.Resource, resource)
		goods[relatedId] = newGood
	}

	sql = "SELECT type,related_id,type_id FROM relation_type WHERE related_id in (" + strings.Join(goodIds, ",") + ") AND delete_at = 0"
	var typeIds []string
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var typeNum int
		var relatedId int64
		var typeId int64
		res.Scan(&typeNum, &relatedId, &typeId)
		var typeObj RelationType
		typeObj.Name = RelationMap[typeNum]
		typeObj.RelatedId = relatedId
		typeObj.TypeId = typeId
		typeIds = append(typeIds, strconv.FormatInt(typeId, 10))
		newGood := goods[relatedId]
		newGood.Type = append(newGood.Type, typeObj)
		goods[relatedId] = newGood
	}

	typeMapping := make(map[int64]string)
	sql = "SELECT id,name FROM type WHERE delete_at = 0 AND id in (" + strings.Join(typeIds, ",") + ")"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var id int64
		var name string
		res.Scan(&id, &name)
		typeMapping[id] = name
	}
	for i := range goods {
		for j := range goods[i].Type {
			goods[i].Type[j].Value = typeMapping[goods[i].Type[j].TypeId]
		}
	}

	users := make(map[int64]User)
	sql = "SELECT id,nickname,mobile,active_at FROM user WHERE id in (" + strings.Join(userIds, ",") + ")"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var id int64
		var nickname string
		var mobile int
		var activeAt int
		res.Scan(&id, &nickname, &mobile, &activeAt)
		var user User
		user.Id = id
		user.Nickname = nickname
		user.Mobile = mobile
		user.ActiveAt = activeAt
		users[id] = user
	}

	sql = "SELECT resource,related_id FROM resource WHERE related_id in (" + strings.Join(userIds, ",") + ") AND related_type = 100 AND delete_at = 0 AND status = 100"
	res, _ = base.Conf.Mysql.Query(sql)
	for res.Next() {
		var resource string
		var relatedId int64
		newUser := users[relatedId]
		newUser.Avatar = resource
		users[relatedId] = newUser
	}

	for i := range goods {
		newGood := goods[i]
		newGood.User = users[goods[i].UserId]
		goods[i]= newGood
	}
	return goods
}
