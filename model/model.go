package model

type Article struct {
	Pic     string `json:"pic"`
	Header  string `json:"header"`
	Content string `json:"content"`
	ID      string `json:"id"`
	Catid   string `json:"catid"`
	Views   int64  `gorm:"default:1" json:"views"`
}

type Picture struct {
	Pic []byte `json:"pic"`
}

type Category struct {
	ID       string `json:"id"`
	Parentid string `json:"parentid"`
	Name     string `json:"name"`
}
