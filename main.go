package main

import (
	"log"

	"github.com/TetsuyaXD/evade/app"
)

func main() {
	application := app.NewApplication()
	if err := application.Setup(); err != nil {
		log.Fatal(err)
		return
	}
	application.Run()
}
