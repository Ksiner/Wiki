package model

type database interface {
	SelectArticles() ([]*Article, error)
	SelectArticlesByCatId(string) ([]*Article, error)
	SelectArticle(string) (*Article, error)
	UpdateArticle(*Article) error
	InsertArticle(*Article) error
	InsertCategory(*Category) error
}

type Article struct {
	Pic     *Picture
	Header  string
	Content string
	Id      string
	Cat     *Category
}

type Picture struct {
	Name string
	Pic  []byte
}

type Category struct {
	Id   string
	Name string
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

func (m *Model) GetArticlesByCatId(catId string) ([]*Article, error) {
	return m.db.SelectArticlesByCatId(catId)
}

func (m *Model) Getrticle(artId string) (*Article, error) {
	return m.db.SelectArticle(artId)
}

func (m *Model) EditArticle(article *Article) error {
	return m.db.UpdateArticle(article)
}

func (m *Model) AddArticle(article *Article) error {
	return m.db.InsertArticle(article)
}

func (m *Model) AddCategory(category *Category) error {
	return m.db.InsertCategory(category)
}
