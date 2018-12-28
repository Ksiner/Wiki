package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Ksiner/Wiki/model"
	"github.com/Ksiner/Wiki/services/db"
	"github.com/Ksiner/Wiki/services/tree"
	"github.com/gorilla/mux"
)

type Config struct {
	Assets string
}

func Start(db db.DataBase, l net.Listener, cfg Config, cancelFunc context.CancelFunc) {
	server := http.Server{
		Handler:        routing(db, cfg),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	go func() {
		defer cancelFunc()
		server.Serve(l)
	}()
}

func routing(db db.DataBase, cfg Config) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, cfg.Assets+"index.html") })
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, cfg.Assets+"favicon.ico") })
	r.Handle("/{filegroup:(?:css|scripts)}/{filename}", getFiles(cfg))
	r.Handle("/{category}/{id}", getArticlesByCatID(db)).Methods("GET")
	r.Handle("/{category}/article/{id}", getArticle(db)).Methods("GET")
	r.HandleFunc("/{category}/article/{path}", getPicture).Methods("GET")
	r.Handle("/{category}/article/{id}/{path}", sendPicture(db)).Methods("POST")
	r.Handle("/{category}/article", editArticle(db, cfg)).Methods("POST")
	r.Handle("/{category}/article/create", createArticle(db))
	r.Handle("/{category}/create", createCategory(db)).Methods("POST")
	r.Handle("/init", getArticles(db, cfg))
	r.Handle("/catTree", getTree(db))
	r.Handle("/login", loginUser(db))
	r.Handle("/logout", logOutUser(db))
	r.Handle("/reg", register(db))
	r.Handle("/token", createToken(db))
	r.Handle("/deleteArts", deleteArts(db))
	r.Handle("/deleteCats", deleteCats(db))
	return r
}

func getFiles(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.ServeFile(w, r, cfg.Assets+vars["filegroup"]+"/"+vars["filename"])
	})
}

func deleteArts(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arts := make([]*model.Article, 0)
		if err := deserialize(r, &arts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.DeleteArticles(arts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func deleteCats(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cats := make([]*model.Category, 0)
		if err := deserialize(r, &cats); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, cat := range cats {
			fmt.Printf("id: %v \n", cat.ID)
			fmt.Printf("Name: %v \n", cat.Name)
			fmt.Printf("Parentid: %v \n", cat.Parentid)
		}
		fmt.Print(cats)
		if err := db.DeleteCaregories(cats); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func createToken(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := db.CreateToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = serialazeAndSend(w, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func loginUser(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t := r.FormValue("token"); t != "" {
			if err := db.AddToken(&model.Token{Token: t}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		user := model.User{}
		err := deserialize(r, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := db.AuthUser(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := serialazeAndSend(w, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func logOutUser(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := model.Token{}
		err := deserialize(r, &token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.LogOutUser(&token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte("true")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func register(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := model.User{}
		if err := deserialize(r, &user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := db.RegisterUser(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := serialazeAndSend(w, token); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func getTree(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, err := tree.New(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = serialazeAndSend(w, t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func getArticles(db db.DataBase, cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articles, err := db.SelectArticles()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		artsWithPics := make([]*model.ArticleWithPic, 0)
		for _, article := range articles {
			withPic := &model.ArticleWithPic{
				ID:      article.ID,
				Catid:   article.Catid,
				Header:  article.Header,
				Content: article.Content,
				Pic:     article.Pic,
				Views:   article.Views,
			}
			pictureBytes, _ := ioutil.ReadFile(cfg.Assets + "picture/" + withPic.Pic)
			if pictureBytes != nil {
				withPic.Picture = pictureBytes
			}
			artsWithPics = append(artsWithPics, withPic)
		}
		err = serialazeAndSend(w, artsWithPics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func getArticlesByCatID(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		articles, err := db.SelectArticlesByCatId(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = serialazeAndSend(w, articles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func getArticle(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		article, err := db.SelectArticle(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		article.Views++
		db.UpdateArticle(article, true)
		err = serialazeAndSend(w, article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func getPicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pictureBytes, err := ioutil.ReadFile("/picture/" + vars["path"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	picture := model.Picture{Pic: pictureBytes}
	serialazeAndSend(w, picture)
}

func sendPicture(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if _, err := os.Stat("/picture/" + vars["path"]); err == nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		picture := model.Picture{}
		deserialize(r, &picture)
		if err := ioutil.WriteFile(vars["path"], picture.Pic, 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := db.UpdateArticlePic(vars["id"], vars["path"]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func editArticle(db db.DataBase, cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.Header.Get("token"))
		articleWithPic := model.ArticleWithPic{}
		err := deserialize(r, &articleWithPic)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := os.Stat(cfg.Assets + "picture/" + articleWithPic.Pic); os.IsNotExist(err) {
			if err := ioutil.WriteFile(cfg.Assets+"picture/"+articleWithPic.Pic, articleWithPic.Picture, 0644); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}
		article := model.Article{
			ID:      articleWithPic.ID,
			Catid:   articleWithPic.Catid,
			Header:  articleWithPic.Header,
			Content: articleWithPic.Content,
			Pic:     articleWithPic.Pic,
			Views:   articleWithPic.Views,
		}
		err = db.UpdateArticle(&article, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
}

func createArticle(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.Header.Get("token"))
		vars := mux.Vars(r)
		article := model.Article{}
		bufferObj := model.BufferArt{}
		err := deserialize(r, &bufferObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		article.Header = bufferObj.Art
		article.Catid = vars["category"]
		err = db.InsertArticle(&article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := serialazeAndSend(w, article); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func createCategory(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print(r.Header.Get("token"))
		vars := mux.Vars(r)
		category := model.Category{}
		bufferObj := model.BufferCat{}
		err := deserialize(r, &bufferObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		category.Name = bufferObj.Cat
		if vars["category"] != "null" {
			category.Parentid = vars["category"]
		}
		err = db.InsertCategory(&category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = serialazeAndSend(w, &category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func deserialize(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func serialazeAndSend(w http.ResponseWriter, v interface{}) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := w.Write(res); err != nil {
		return err
	}
	return nil
}
