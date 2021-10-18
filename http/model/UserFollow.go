package model

type UserFollow struct {
	Id       int64 `json:"id"`
	UserId   int64 `json:"user_id"`
	FollowId int64 `json:"follow_id"`
	CreateAt int   `json:"create_at"`
	UpdateAt int   `json:"update_at"`
	DeleteAt int   `json:"delete_at"`
}
