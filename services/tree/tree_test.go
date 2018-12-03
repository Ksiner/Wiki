package tree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Ksiner/Wiki/model"

	"github.com/Ksiner/Wiki/services/db"
)

func GetConfigs() (*db.DbConnMysql, error) {
	var cfg db.Config
	byteCfg, err := ioutil.ReadFile("/home/ksiner/go/src/github.com/Ksiner/Wiki/services/db/testingconf.json")
	if err != nil {
		fmt.Printf("Error in reading tests config json file: %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal(byteCfg, &cfg)
	if err != nil {
		fmt.Printf("Error in parsing tests config json file: %v", err.Error())
		return nil, err
	}
	return &db.DbConnMysql{Cfg: cfg}, nil
}

func TestNewTree(t *testing.T) {
	dbc, err := GetConfigs()
	if err != nil {
		t.Errorf("Error during reading config file! %v", err.Error())
	}
	trees, err := New(dbc)
	if err != nil {
		t.Errorf("Error during making tree! %v", err.Error())
	}
	for _, tree := range trees {
		openTree(0, tree)
	}
}

func openTree(level int, tree *model.CatTree) {
	fmt.Printf("\nLevel: %v, Value: %v \n", level, tree.Cat.Name)
	for _, child := range tree.Childs {
		openTree(level+1, child)
	}
}
