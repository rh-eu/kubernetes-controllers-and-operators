package main

import (
	"log"

	"github.com/rh-eu/kubernetes-controllers-and-operators/pkg/app"
)

func main() {
	app := app.NewApp()

	log.Printf("Starting app ... %+v", app)

	app.Run()
}
