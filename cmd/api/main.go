package main

import (
	"log"

	"github.com/canxium/supply-information/config"
	"github.com/canxium/supply-information/server"
)

func main() {
	log.Println("Starting canxium supply api server")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	s := server.NewServer(cfg)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
