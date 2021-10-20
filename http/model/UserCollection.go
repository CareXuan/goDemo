package model

type UserCollection struct {
	Id       int64 `json:"-"`
	UserId   int64 `json:"user_id"`
	GoodId   int64 `json:"good_id"`
	CreateAt int   `json:"-"`
	UpdateAt int   `json:"-"`
	DeleteAt int   `json:"-"`
}
