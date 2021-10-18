package model

type User struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Mobile   int    `json:"mobile"`
	CreateAt int `json:"create_at"`
	UpdateAt int `json:"update_at"`
	DeleteAt int `json:"delete_at"`
}
