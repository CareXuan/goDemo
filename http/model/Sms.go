package model

type Sms struct {
	Id     int64 `json:"id"`
	Type   int   `json:"type"`
	Mobile int   `json:"mobile"`
	Code   int   `json:"code"`
	Status int   `json:"status"`
}
