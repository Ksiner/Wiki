package daemon

import (
	"context"
	"fmt"
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

func Run(cfg *Config) (context.Context, error) {
	db := db.NewMySql(cfg.DBConfig)
	l, err := net.Listen(cfg.Network, cfg.Address)
	if err != nil {
		fmt.Printf("Error in daemon! %v", err.Error())
		return nil, err
	}
	context, cancelFunc := context.WithCancel(context.Background())
	ui.Start(db, l, cfg.UIConfig, cancelFunc)
	return context, nil
}
