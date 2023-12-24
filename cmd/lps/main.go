package main

import (
	"log"

	"lps/internal/app"
	"lps/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
