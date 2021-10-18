package model

type UserCollection struct {
	Id       int64 `json:"id"`
	UserId   int64 `json:"user_id"`
	GoodId   int64 `json:"good_id"`
	CreateAt int   `json:"create_at"`
	UpdateAt int   `json:"update_at"`
	DeleteAt int   `json:"delete_at"`
}
