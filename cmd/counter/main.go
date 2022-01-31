package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	app := createApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
