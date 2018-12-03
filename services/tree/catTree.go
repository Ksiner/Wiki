package tree

import (
	"github.com/Ksiner/Wiki/model"
	"github.com/Ksiner/Wiki/services/db"
)

func New(db db.DataBase) ([]*model.CatTree, error) {
	cats, err := db.SelectCategories()
	if err != nil {
		return nil, err
	}
	arts, err := db.SelectArticles()
	if err != nil {
		return nil, err
	}
	trees := make([]*model.CatTree, 0)
	for _, cat := range cats {
		if cat.Parentid == "" {
			tree, err := makeTree(cat, nil, cats, arts)
			if err != nil {
				return nil, err
			}
			trees = append(trees, tree)
		}
	}
	if err != nil {
		return nil, err
	}
	return trees, nil
}

func makeTree(currCat *model.Category, parent *model.Category, cats []*model.Category, arts []*model.Article) (*model.CatTree, error) {
	tree := &model.CatTree{Cat: currCat, Parent: parent}
	tree.Articles = findArticles(tree.Cat, arts)
	childs := findChilds(currCat, cats)
	for _, child := range childs {
		childTree, err := makeTree(child, currCat, cats, arts)
		if err != nil {
			return nil, err
		}
		tree.Childs = append(tree.Childs, childTree)
	}
	return tree, nil
}

func findArticles(currCat *model.Category, arts []*model.Article) []*model.Article {
	res := make([]*model.Article, 0)
	for _, art := range arts {
		if art.Catid == currCat.ID {
			res = append(res, art)
		}
	}
	return res
}

func findChilds(currCat *model.Category, cats []*model.Category) []*model.Category {
	childs := make([]*model.Category, 0)
	for _, child := range cats {
		if child.Parentid == currCat.ID {
			childs = append(childs, child)
		}
	}
	return childs
}
