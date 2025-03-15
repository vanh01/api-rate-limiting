package main

import (
	"log"

	"github.com/vanh01/api-rate-limiting/config"
	"github.com/vanh01/api-rate-limiting/internal/app"
)

func main() {
	// load config from file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config :%s", err.Error())
	}
	config.Instance = cfg

	// run application
	app.Run()
}
