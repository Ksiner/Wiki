package daemon

import (
	"log"
	"net"

	"github.com/Ksiner/Wiki/ctrl/ui"
	"github.com/Ksiner/Wiki/services/db"
)

type Config struct {
	Network  string
	Address  string
	DBConfig db.Config
	UIConfig ui.Config
}

func Run(cfg Config) error {
	db := db.NewMySql(cfg.DBConfig)
	l, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		log.Panicf("Error in daemon! %v", err.Error())
		return err
	}
	ui.Start(db, l, cfg.UIConfig)
	return nil
}
