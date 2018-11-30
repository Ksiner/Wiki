package model

type Article struct {
	Pic     string
	Header  string
	Content string
	ID      string
	Catid   string
	Views   int64 `gorm:"default:1"`
}

type Picture struct {
	Pic []byte `json:"pic"`
}

type Category struct {
	ID       string
	Parentid string
	Name     string
}
