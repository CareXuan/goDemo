package model

type Sms struct {
	Id       int64 `json:"-"`
	Type     int   `json:"-"`
	Mobile   int   `json:"mobile"`
	Code     int   `json:"code"`
	Status   int   `json:"status"`
	UsedAt   int   `json:"used_at"`
	CreateAt int   `json:"-"`
	UpdateAt int   `json:"-"`
	DeleteAt int   `json:"-"`
}
