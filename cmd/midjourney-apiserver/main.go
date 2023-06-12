package main

import (
	"log"

	"github.com/hongliang5316/midjourney-apiserver/internal/application"
)

func main() {
	app := application.New()
	if err := app.Run(); err != nil {
		log.Fatalf("Call app.Run failed, err: %+v", err)
	}
}
