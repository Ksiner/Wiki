package model

type Article struct {
	Pic     string `json:"pic"`
	Header  string `json:"header"`
	Content string `json:"content"`
	ID      string `json:"id"`
	Catid   string `json:"catid"`
	Views   int64  `gorm:"default:1" json:"views"`
}

type ArticleWithPic struct {
	Pic     string `json:"pic"`
	Header  string `json:"header"`
	Content string `json:"content"`
	ID      string `json:"id"`
	Catid   string `json:"catid"`
	Views   int64  `json:"views"`
	Picture []byte `json:"picture"`
}

type Picture struct {
	Pic []byte `json:"pic"`
}

type CatTree struct {
	Cat      *Category  `json:"category"`
	Parent   *Category  `json:"parent"`
	Childs   []*CatTree `json:"childs"`
	Articles []*Article `json:"articles"`
}

type Category struct {
	ID       string `json:"id"`
	Parentid string `json:"parentid"`
	Name     string `json:"name"`
}

type User struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Pass  string `json:"pass"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Token struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type BufferArt struct {
	Art string `json:"art"`
}

type BufferCat struct {
	Cat string `json:"cat"`
}
