package db

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"

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

func (dbc *DbConnMysql) SelectCategories() ([]*model.Category, error) {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectCategories\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	cats := make([]*model.Category, 0)
	db.Find(&cats)
	if db.Error != nil {
		fmt.Printf("Error in \"db.SelectCategories\" GET DATA: %v", err.Error())
		err = db.Error
		return nil, err
	}
	return cats, nil
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
	category.ID = uuid.Must(uuid.NewV4()).String()
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

func (dbc *DbConnMysql) DeleteArticles(arts []*model.Article) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.DeleteArticles\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()

	db.Delete(&arts)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.DeleteArticles\" DELETE ARTS: %v", err.Error())
		return err
	}
	return nil
}

func (dbc *DbConnMysql) DeleteCaregories(cats []*model.Category) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.DeleteCaregories\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()

	db.Delete(&cats)

	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.DeleteCaregories\" DELETE CATS: %v", err.Error())
		return err
	}
	return nil
}

func (dbc *DbConnMysql) CreateToken() (*model.Token, error) {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.SelectArticles\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	token := model.Token{Token: uuid.Must(uuid.NewV4()).String(), Role: "reader"}
	db.Save(&token)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.RegisterUser\" INSERT AUTH TOKEN: %v", err.Error())
		return nil, err
	}
	return &token, nil
}

func (dbc *DbConnMysql) AddToken(token *model.Token) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.AddToken\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	db.Save(&token)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.AddToken\" INSERT AUTH TOKEN: %v", err.Error())
		return err
	}
	return nil
}

func (dbc *DbConnMysql) CheckToken(token *model.Token) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.CheckToken\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	checkToken := model.Token{}
	db.Where("token=?", token.Token).First(&checkToken)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.CheckToken\" INSERT AUTH TOKEN: %v", err.Error())
		return err
	}
	if checkToken.Token == "" {
		return errors.New("No such token!")
	}
	return nil
}

func (dbc *DbConnMysql) AuthUser(user *model.User) (*model.Token, error) {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.AuthUser\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	userCheck := model.User{}
	db.Where("(login = ? or email = ?) and pass = ?", user.Login, user.Email, user.Pass).First(&userCheck)
	if userCheck.ID == "" {
		return nil, errors.New(user.Login + " no such user!")
	}
	token := model.Token{Token: uuid.Must(uuid.NewV4()).String(), Role: userCheck.Role}
	db.Save(&token)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.AuthUser\" INSERT AUTH TOKEN: %v", err.Error())
		return nil, err
	}
	return &token, nil
}

func (dbc *DbConnMysql) LogOutUser(token *model.Token) error {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.LogOutUser\" OPEN DB: %v", err.Error())
		return err
	}
	defer db.Close()
	db.Delete(&token)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.LogOutUser\" DELETE AUTH TOKEN: %v", err.Error())
		return err
	}
	return nil
}

func (dbc *DbConnMysql) RegisterUser(user *model.User) (*model.Token, error) {
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.LogOutUser\" OPEN DB: %v", err.Error())
		return nil, err
	}
	defer db.Close()
	user.ID = uuid.Must(uuid.NewV4()).String()
	db.Save(&user)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.RegisterUser\" INSERT USER: %v", err.Error())
		return nil, err
	}
	token := model.Token{Token: uuid.Must(uuid.NewV4()).String(), Role: "reader"}
	db.Save(&token)
	if db.Error != nil {
		err = db.Error
		db.Error = nil
		fmt.Printf("Error in \"db.RegisterUser\" INSERT AUTH TOKEN: %v", err.Error())
		return nil, err
	}
	return &token, nil
}
