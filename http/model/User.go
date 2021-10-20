package model

import (
	"fmt"
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
	Avatar   string `json:"avatar"`
	CreateAt int    `json:"-"`
	UpdateAt int    `json:"-"`
	DeleteAt int    `json:"-"`
}

var userObj User

func GetOneUserByMobile(mobile string) User {
	sql := "SELECT user_id,nickname,mobile,token FROM user WHERE mobile = ? and delete_at = 0"
	err := base.Conf.Mysql.QueryRow(sql, mobile).Scan(
		&userObj.UserId,
		&userObj.Nickname,
		&userObj.Mobile,
		&userObj.Token,
	)
	if err != nil {
		fmt.Print(err)
	}
	sql = "SELECT resource FROM resource WHERE related_type = 100 and related_id = ?"
	_ = base.Conf.Mysql.QueryRow(sql, userObj.UserId).Scan(&userObj.Avatar)
	return userObj
}

func CreateOneUser(mobile string) User {
	sql := "INSERT INTO user(user_id,nickname,password,mobile,create_at,update_at) VALUES(?,?,?,?,?,?)"
	_ = base.Conf.Mysql.QueryRow(sql, strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(rand.Intn(500)), "忙猫_"+common.RandString(6), "", mobile, time.Now().Unix(), time.Now().Unix())
	return GetOneUserByMobile(mobile)
}
