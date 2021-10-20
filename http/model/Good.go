package model

type Good struct {
	Id       int64   `json:"-"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Intro    string  `json:"intro"`
	Type     int     `json:"type"`
	UserId   int64   `json:"user_id"`
	Status   int     `json:"status"`
	CreateAt int     `json:"-"`
	UpdateAt int     `json:"-"`
	DeleteAt int     `json:"-"`
}
