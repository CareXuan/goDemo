package model

type RelationType struct {
	Id        int64  `json:"-"`
	Type      int    `json:"-"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	RelatedId int64  `json:"-"`
	TypeId    int64  `json:"-"`
	CreateAt  int    `json:"-"`
	UpdateAt  int    `json:"-"`
	DeleteAt  int    `json:"-"`
}

var RelationMap = map[int]string{
	100: "成色",
}
