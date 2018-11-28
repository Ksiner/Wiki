package model

type database interface {
	SelectArticles() ([]*Article, error)
	SelectArticlesByCatId(string) ([]*Article, error)
	SelectArticle(string) (*Article, error)
	UpdateArticle(*Article, bool) error
	InsertArticle(*Article) error
	InsertCategory(*Category) error
}

type Article struct {
	Pic     string
	Header  string
	Content string
	ID      string
	Catid   string
	Views   int64 `gorm:"default:1"`
}

type Category struct {
	ID       string
	Parentid string
	Name     string
}

type Model struct {
	db database
}

func New(db database) *Model {
	return &Model{db: db}
}

func (m *Model) GetArticles() ([]*Article, error) {
	return m.db.SelectArticles()
}

func (m *Model) GetArticlesByCatId(catID string) ([]*Article, error) {
	return m.db.SelectArticlesByCatId(catID)
}

func (m *Model) Getrticle(artID string) (*Article, error) {
	return m.db.SelectArticle(artID)
}

func (m *Model) EditArticle(article *Article, addView bool) error {
	return m.db.UpdateArticle(article, addView)
}

func (m *Model) AddArticle(article *Article) error {
	return m.db.InsertArticle(article)
}

func (m *Model) AddCategory(category *Category) error {
	return m.db.InsertCategory(category)
}
