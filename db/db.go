package db

import (
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/ksiner/wiki/model"
)

type Config struct {
	dbConnStr string
	dbDialect string
}

type DbConn struct {
	cfg Config
}

func New(cfg Config) *DbConn {
	return &DbConn{cfg: cfg}
}

func (dbc *DbConn) SelectArticles() ([]*model.Article, error) {
	db, err := gorm.Open(dbc.cfg.dbDialect, dbc.cfg.dbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectArticles\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	articles := make([]*model.Article, 0)
	db.Find(articles).Order("views desc")
	if db.Error != nil {
		fmt.Printf("Error in \"db.SelectArticles\" GET DATA: %v", err.Error())
		err = db.Error
		return nil, err
	}
	return articles, nil
}

func (dbc *DbConn) SelectArticlesByCatId(catID string) ([]*model.Article, error) {
	db, err := gorm.Open(dbc.cfg.dbDialect, dbc.cfg.dbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectArticlesByCatId\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	articles := make([]*model.Article, 0)
	db.Find(articles).Where("catid = ?", catID).Order("views desc")
	if db.Error != nil {
		fmt.Printf("Error in \"db.SelectArticlesByCatId\" GET DATA: %v", err.Error())
		err = db.Error
		db.Error = nil
		return nil, err
	}
	return articles, nil
}

func (dbc *DbConn) SelectArticle(artID string) (*model.Article, error) {
	db, err := gorm.Open(dbc.cfg.dbDialect, dbc.cfg.dbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectArticle\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	article := &model.Article{}
	db.First(article).Where("id = ?", artID).Order("views desc")
	if db.Error != nil {
		fmt.Printf("Error in \"db.SelectArticle\" GET DATA: %v", err.Error())
		err = db.Error
		db.Error = nil
		return nil, err
	}
	return article, nil
}

func (dbc *DbConn) InsertArticle(article *model.Article) error {
	db, err := gorm.Open(dbc.cfg.dbDialect, dbc.cfg.dbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.InsertArticle\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	for db.NewRecord(article) != true {
		article.ID = uuid.Must(uuid.NewV4()).String()
	}
	db.Create(article)
	if db.Error != nil {
		fmt.Printf("Error in \"db.InsertArticle\" INSERT DB: %v", err.Error())
		err = db.Error
		db.Error = nil
		return err
	}
	return nil
}

func (dbc *DbConn) InsertCategory(category *model.Category) error {
	db, err := gorm.Open(dbc.cfg.dbDialect, dbc.cfg.dbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.InsertCategory\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	for db.NewRecord(category) != true {
		category.ID = uuid.Must(uuid.NewV4()).String()
	}
	db.Create(category)
	if db.Error != nil {
		fmt.Printf("Error in \"db.InsertCategory\" INSERT DB: %v", err.Error())
		err = db.Error
		db.Error = nil
		return err
	}
	return nil
}

func (dbc *DbConn) UpdateArticle(article *model.Article, addView bool) error {
	db, err := gorm.Open(dbc.cfg.dbDialect, dbc.cfg.dbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.UpdateArticle\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	if addView {
		db.Model(article).Update("views", gorm.Expr("views+1"))
		if db.Error != nil {
			fmt.Printf("Error in \"db.UpdateArticle\" UPDATE VIEWS DB: %v", err.Error())
			err = db.Error
			db.Error = nil
			return err
		}
		return nil
	}
	db.Save(article)
	if db.Error != nil {
		fmt.Printf("Error in \"db.UpdateArticle\" UPDATE ALL DB: %v", err.Error())
		err = db.Error
		db.Error = nil
		return err
	}
	return nil
}
