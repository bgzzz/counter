package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	app := createCtlApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
