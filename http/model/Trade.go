package model

type Trade struct {
	Id          int64   `json:"-"`
	UserId      int64   `json:"user_id"`
	GoodId      int64   `json:"good_id"`
	Price       float64 `json:"price"`
	FinishPrice float64 `json:"finish_price"`
	Status      int     `json:"status"`
	CreateAt    int     `json:"-"`
	UpdateAt    int     `json:"-"`
	DeleteAt    int     `json:"-"`
}
