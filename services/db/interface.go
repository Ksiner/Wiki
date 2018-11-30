package db

import "github.com/Ksiner/Wiki/model"

type DataBase interface {
	SelectArticles() ([]*model.Article, error)
	SelectArticlesByCatId(string) ([]*model.Article, error)
	SelectArticle(string) (*model.Article, error)
	UpdateArticlePic(string, string) error
	UpdateArticle(*model.Article, bool) error
	InsertArticle(*model.Article) error
	InsertCategory(*model.Category) error
}
