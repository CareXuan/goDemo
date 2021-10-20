package model

type Resource struct {
	Id          int64  `json:"-"`
	RelatedType int64  `json:"-"`
	RelatedId   int64  `json:"-"`
	Resource    string `json:"resource"`
	Status      int    `json:"-"`
	CreateAt    int    `json:"-"`
	UpdateAt    int    `json:"-"`
	DeleteAt    int    `json:"-"`
}
