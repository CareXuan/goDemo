package model

import (
	"math/rand"
	"mouse/base"
	"mouse/common"
	"strconv"
	"time"
)

type User struct {
	Id       int64  `json:"-"`
	UserId   int64  `json:"-"`
	Nickname string `json:"nickname"`
	Password string `json:"-"`
	Mobile   int    `json:"mobile"`
	Token    string `json:"-"`
	CreateAt int    `json:"-"`
	UpdateAt int    `json:"-"`
	DeleteAt int    `json:"-"`
}

var userObj User

func GetOneUserByMobile(mobile string) User {
	sql := "SELECT user_id,nickname,mobile,token FROM user WHERE mobile = ? and delete_at = 0"
	res, _ := base.Conf.Mysql.Query(sql, mobile)
	for res.Next() {
		var userId int64
		var nickname string
		var mobile int
		var token string
		res.Scan(&userId, &nickname, &mobile, &token)
		userObj.Nickname = nickname
		userObj.UserId = userId
		userObj.Mobile = mobile
		userObj.Token = token
	}
	return userObj
}

func CreateOneUser(mobile string) User {
	sql := "INSERT INTO user(user_id,nickname,password,mobile,create_at,update_at) VALUES(?,?,?,?,?,?)"
	_ = base.Conf.Mysql.QueryRow(sql, strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(rand.Intn(500)), "忙猫_"+common.RandString(6), "", mobile, time.Now().Unix(), time.Now().Unix())
	return GetOneUserByMobile(mobile)
}
