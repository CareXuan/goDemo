package model

type Good struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Intro  string  `json:"intro"`
	Type   int     `json:"type"`
	UserId int64   `json:"user_id"`
	Status int     `json:"status"`
	CreateAt int `json:"create_at"`
	UpdateAt int `json:"update_at"`
	DeleteAt int `json:"delete_at"`
}
