package model

type UserFollow struct {
	Id       int64 `json:"-"`
	UserId   int64 `json:"user_id"`
	FollowId int64 `json:"follow_id"`
	CreateAt int   `json:"-"`
	UpdateAt int   `json:"-"`
	DeleteAt int   `json:"-"`
}
