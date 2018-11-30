package db

import (
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/Ksiner/Wiki/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Config struct {
	DbConnStr string
	DbDialect string
}

type DbConnMysql struct {
	Cfg Config
}

func NewMySql(cfg Config) *DbConnMysql {
	return &DbConnMysql{Cfg: cfg}
}

func (dbc *DbConnMysql) SelectArticles() ([]*model.Article, error) {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectArticles\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	articles := make([]*model.Article, 0)
	db.Find(&articles).Order("views desc")
	if db.Error != nil {
		fmt.Printf("Error in \"db.SelectArticles\" GET DATA: %v", err.Error())
		err = db.Error
		return nil, err
	}
	return articles, nil
}

func (dbc *DbConnMysql) SelectArticlesByCatId(catID string) ([]*model.Article, error) {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectArticlesByCatId\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	articles := make([]*model.Article, 0)
	db.Find(&articles).Where("catid = ?", catID).Order("views desc")
	if db.Error != nil {
		fmt.Printf("Error in \"db.SelectArticlesByCatId\" GET DATA: %v", err.Error())
		err = db.Error
		db.Error = nil
		return nil, err
	}
	return articles, nil
}

func (dbc *DbConnMysql) SelectArticle(artID string) (*model.Article, error) {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectArticle\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	article := &model.Article{}
	db.Where("id = ?", artID).First(&article)
	if db.Error != nil {
		fmt.Printf("Error in \"db.SelectArticle\" GET DATA: %v", err.Error())
		err = db.Error
		db.Error = nil
		return nil, err
	}
	return article, nil
}

func (dbc *DbConnMysql) InsertArticle(article *model.Article) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.InsertArticle\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	article.ID = uuid.Must(uuid.NewV4()).String()
	db.Save(article)
	if db.Error != nil {
		fmt.Printf("Error in \"db.InsertArticle\" INSERT DB: %v", err.Error())
		err = db.Error
		db.Error = nil
		return err
	}
	return nil
}

func (dbc *DbConnMysql) InsertCategory(category *model.Category) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.InsertCategory\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	db.Save(category)
	if db.Error != nil {
		fmt.Printf("Error in \"db.InsertCategory\" INSERT DB: %v", err.Error())
		err = db.Error
		db.Error = nil
		return err
	}
	return nil
}

func (dbc *DbConnMysql) UpdateArticle(article *model.Article, justViews bool) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.UpdateArticle\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	if justViews {
		db.Model(article).UpdateColumn("views", gorm.Expr("views+1"))
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

func (dbc *DbConnMysql) UpdateArticlePic(artID string, path string) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.UpdateArticlePic\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	db.Table("articles").Where("id = ?", artID).UpdateColumn("pic", gorm.Expr("?", path))
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.UpdateArticlePic\" UPDATE PIC DB: %v", err.Error())
		return err
	}
	return nil
}
