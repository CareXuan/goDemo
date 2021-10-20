package model

type Good struct {
	Id       int64    `json:"-"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Intro    string   `json:"intro"`
	Type     int      `json:"-"`
	UserId   int64    `json:"-"`
	Status   int      `json:"-"`
	Resource []string `json:"resource"`
	CreateAt int      `json:"-"`
	UpdateAt int      `json:"-"`
	DeleteAt int      `json:"-"`
}
