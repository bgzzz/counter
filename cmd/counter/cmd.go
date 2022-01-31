package main

import (
	"github.com/bgzzz/counter/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func createApp() *cli.App {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "log-level \"debug\" (more on the supported levels here: https://github.com/sirupsen/logrus/blob/fdf1618bf7436ec3ee65753a6e2999c335e97221/logrus.go#L25)",
			Value:   "debug",
			EnvVars: []string{"LOG_LEVEL"},
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "--port 8080",
			Value:   8080,
			EnvVars: []string{"PORT"},
		},
	}

	return &cli.App{
		Name:  "counter-server",
		Usage: "counter-server is server side application returning the value of the counter",
		Flags: flags,
		Action: func(c *cli.Context) error {

			log.SetFormatter(&log.JSONFormatter{})
			lvl, err := log.ParseLevel(c.String("log-level"))
			if err != nil {
				log.Errorf("unable to parse log level (%s): %v",
					c.String("log-level"), err)
			}
			log.SetLevel(lvl)

			port := c.Int("port")

			return server.NewServer(port).Run()
		},
	}
}
