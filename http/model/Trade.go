package model

type Trade struct {
	Id          int64   `json:"id"`
	UserId      int64   `json:"user_id"`
	GoodId      int64   `json:"good_id"`
	Price       float64 `json:"price"`
	FinishPrice float64 `json:"finish_price"`
	Status      int     `json:"status"`
	CreateAt int `json:"create_at"`
	UpdateAt int `json:"update_at"`
	DeleteAt int `json:"delete_at"`
}
