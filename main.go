package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Ksiner/Wiki/daemon"
)

func ParseConfigJSON() (*daemon.Config, error) {
	cfgBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	cfg := daemon.Config{}
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func WaitForSignal(ctx context.Context) {
	<-ctx.Done()
	log.Panicf("Server is shutted down!")
}

func main() {
	cfg, err := ParseConfigJSON()
	if err != nil {
		log.Panicf("Error in main! %v", err.Error())
	}
	ctx, err := daemon.Run(cfg)
	if err != nil {
		log.Panicf("Error in main! %v", err.Error())
	}
	WaitForSignal(ctx)
}
