package main

import (
	"golangTestTask/internal/app"
	"golangTestTask/internal/config"
	"log"

	"github.com/caarlos0/env"
)

func main() {
	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal("error running http server ", err)
	}
}
