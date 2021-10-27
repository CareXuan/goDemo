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
	ActiveAt int    `json:"active_at"`
	CreateAt int    `json:"-"`
	UpdateAt int    `json:"-"`
	DeleteAt int    `json:"-"`
}

var userObj User

//
// GetOneUserByMobile
// @Description: 根据手机号获取用户
// @param mobile
// @return User
//
func GetOneUserByMobile(mobile string) User {
	sql := "SELECT id,user_id,nickname,mobile,token FROM user WHERE mobile = ? and delete_at = 0"
	err := base.Conf.Mysql.QueryRow(sql, mobile).Scan(
		&userObj.Id,
		&userObj.UserId,
		&userObj.Nickname,
		&userObj.Mobile,
		&userObj.Token,
	)
	if err != nil {
		fmt.Print(err)
		return userObj
	}
	sql = "SELECT resource FROM resource WHERE related_type = 100 and related_id = ?"
	_ = base.Conf.Mysql.QueryRow(sql, userObj.UserId).Scan(&userObj.Avatar)
	return userObj
}

//
// GetOneUserByUserId
// @Description: 根据user_id获取用户
// @param userId
// @return User
//
func GetOneUserByUserId(userId string) User {
	sql := "SELECT id,user_id,nickname,mobile,token FROM user WHERE user_id = ? and delete_at = 0"
	err := base.Conf.Mysql.QueryRow(sql, userId).Scan(
		&userObj.Id,
		&userObj.UserId,
		&userObj.Nickname,
		&userObj.Mobile,
		&userObj.Token,
	)
	if err != nil {
		fmt.Print(err)
		return userObj
	}
	sql = "SELECT resource FROM resource WHERE related_type = 100 and related_id = ?"
	_ = base.Conf.Mysql.QueryRow(sql, userObj.UserId).Scan(&userObj.Avatar)
	return userObj
}

//
// CreateOneUser
// @Description: 创建新用户
// @param mobile
// @return User
//
func CreateOneUser(mobile string) User {
	sql := "INSERT INTO user(user_id,nickname,password,mobile,create_at,update_at) VALUES(?,?,?,?,?,?)"
	_ = base.Conf.Mysql.QueryRow(sql, strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(rand.Intn(500)), "忙猫_"+common.RandString(6), "", mobile, time.Now().Unix(), time.Now().Unix())
	return GetOneUserByMobile(mobile)
}
