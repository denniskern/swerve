package main

import (
	"github.com/hammi85/swerve/src/app"
)

func main() {
	// new application
	application := app.NewApplication()
	// read configuration
	application.Setup()
	// run the server
	application.Run()
}
