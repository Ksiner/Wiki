package ui

import (
	"context"
	"encoding/json"
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
	r.Handle("/{filegroup:(?:css|scripts)}/{filename}", getFiles(cfg))
	r.Handle("/{category}/{id}", getArticlesByCatID(db))
	r.Handle("/{category}/article/{id}", getArticle(db)).Methods("GET")
	r.HandleFunc("/{category}/article/{path}", getPicture).Methods("GET")
	r.Handle("/{category}/article/{id}/{path}", sendPicture(db)).Methods("POST")
	r.Handle("/{category}/article", editArticle(db)).Methods("POST")
	r.Handle("/{category}/article/create", createArticle(db))
	r.Handle("/{category}/create", createCategory(db))
	r.Handle("/init", getArticles(db))
	r.Handle("/catTree", getTree(db))
	return r
}

func getFiles(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.ServeFile(w, r, cfg.Assets+vars["filegroup"]+"/"+vars["filename"])
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

func getArticles(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articles, err := db.SelectArticles()
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

func editArticle(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		article := model.Article{}
		err := deserialize(r, &article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		article := model.Article{}
		err := deserialize(r, &article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.InsertArticle(&article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write([]byte(article.ID)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func createCategory(db db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		category := model.Category{}
		err := deserialize(r, &category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.InsertCategory(&category)
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
