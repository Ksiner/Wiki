package db

import "github.com/Ksiner/Wiki/model"

type DataBase interface {
	SelectArticles() ([]*model.Article, error)
	SelectArticlesByCatId(string) ([]*model.Article, error)
	SelectArticle(string) (*model.Article, error)
	SelectCategories() ([]*model.Category, error)
	UpdateArticlePic(string, string) error
	UpdateArticle(*model.Article, bool) error
	InsertArticle(*model.Article) error
	InsertCategory(*model.Category) error
	AuthUser(*model.User) (*model.Token, error)
	LogOutUser(*model.Token) error
	RegisterUser(*model.User) (*model.Token, error)
	AddToken(*model.Token) error
	CheckToken(*model.Token) error
	CreateToken() (*model.Token, error)
}
