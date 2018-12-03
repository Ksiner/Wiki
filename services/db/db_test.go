package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/Ksiner/Wiki/model"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type User struct {
	ID   int
	Name string
}

func TestSelectArticles(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error in getting tests config:%v", err.Error())
	}
	if articles, err := dbc.SelectArticles(); err != nil {
		t.Error(err.Error())
	} else {
		if cap(articles) == 0 {
			t.Error("Cap of articles is 0, expected >0")
		}
	}
}

func TestSelectArticleOnArticles(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error in getting tests config:%v", err.Error())
	}
	if articles, err := dbc.SelectArticles(); err != nil {
		t.Error(err.Error())
	} else {
		if cap(articles) == 0 {
			t.Error("Cap of articles is 0, expected >0")
		}
		for _, art := range articles {
			selArt, err := dbc.SelectArticle(art.ID)
			if err != nil {
				t.Error(err.Error())
			}
			if !reflect.DeepEqual(selArt, art) {
				fmt.Printf("selArt: ID = %v, Pic = %v, Header= %v, Content = %v, Cat = %v \n", selArt.ID, selArt.Pic, selArt.Header, selArt.Content, selArt.Catid)
				fmt.Printf("Art: ID = %v, Pic = %v, Header= %v, Content = %v, Cat = %v \n", art.ID, art.Pic, art.Header, art.Content, art.Catid)
				t.Errorf("Articles are not equal: selArt is %v, and Art is %v", selArt, art)
			}
		}
	}
}

func TestSelectArticleNone(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error in getting tests config:%v", err.Error())
	}
	article, err := dbc.SelectArticle("")
	if err != nil {
		t.Error(err.Error())
	}
	if article.ID != "" {
		t.Error(article)
	}
}

func TestToInsertSomething(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error in getting tests config:%v", err.Error())
	}
	db, err := gorm.Open(dbc.Cfg.DbDialect, dbc.Cfg.DbConnStr)
	if err != nil {
		fmt.Printf("Error in \"db.InsertArticle\" OPEN DB: %v", err.Error())
		t.Error(err.Error())
	}
	defer db.Close()
	art := model.Article{}
	article := model.Article{ID: uuid.Must(uuid.NewV4()).String(), Catid: "12312", Content: "some testing content12", Header: "12312", Pic: "some testing pic2", Views: 0}
	db.FirstOrCreate(&art, article)
}

func TestInsertSimpleArticle(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error in getting tests config:%v", err.Error())
	}
	article := &model.Article{ID: uuid.Must(uuid.NewV4()).String(), Catid: "some testing category", Content: "some testing content", Header: "some testing header", Pic: "some testing pic", Views: 1}
	err = dbc.InsertArticle(article)
	if err != nil {
		t.Error(err.Error())
	}
	selectedArticle, err := dbc.SelectArticle(article.ID)
	if err != nil {
		t.Errorf("Error on select inserted article: %v", err.Error())
	}
	if !reflect.DeepEqual(selectedArticle, article) {
		t.Errorf("Inserted and selected row are not equals!\n inserted: %v \n selected %v \n", article, selectedArticle)
	}
}

func TestUpdateArticleJsutViews(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error in getting tests config:%v", err.Error())
	}
	article := model.Article{ID: "44c725d3-2a4c-436b-928d-d36d1addf265", Pic: "some testing pic2", Header: "12312", Content: "some testing content12", Catid: "12312"}
	err = dbc.UpdateArticle(&article, true)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestUpdateArticleFully(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error in getting tests config:%v", err.Error())
	}
	article := model.Article{ID: "44c725d3-2a4c-436b-928d-d36d1addf265", Pic: "some tesic2", Header: "112", Content: "some ntent12", Views: 1, Catid: "1"}
	err = dbc.UpdateArticle(&article, false)
	if err != nil {
		t.Error(err.Error())
	}
}

func GetConfigs() (*DbConnMysql, error) {
	var cfg Config
	byteCfg, err := ioutil.ReadFile("testingconf.json")
	if err != nil {
		fmt.Printf("Error in reading tests config json file: %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal(byteCfg, &cfg)
	if err != nil {
		fmt.Printf("Error in parsing tests config json file: %v", err.Error())
		return nil, err
	}
	return &DbConnMysql{Cfg: cfg}, nil
}
