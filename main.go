package main

import (
	"log"

	"github.com/munye/cardioprieto-api/app"
	"github.com/munye/cardioprieto-api/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := config.NewConfig()
	app.ConfigAndRunApp(config)
}
