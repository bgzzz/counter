package main

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	cl "github.com/bgzzz/counter/pkg/client"
	"github.com/bgzzz/counter/pkg/model"
)

const protocolSchemaTemplate = "http://%s"

func createCtlApp() *cli.App {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "log-level \"debug\" (more on the supported levels here: https://github.com/sirupsen/logrus/blob/fdf1618bf7436ec3ee65753a6e2999c335e97221/logrus.go#L25)",
			Value:   "debug",
			EnvVars: []string{"LOG_LEVEL"},
		},
		&cli.StringFlag{
			Name:    "host",
			Usage:   "--host localhost:8080",
			Value:   "localhost:8080",
			EnvVars: []string{"host"},
		},
	}

	return &cli.App{
		Name:  "counterctl",
		Usage: "counterctl is a client side application calling the counter server",
		Flags: flags,
		Commands: []*cli.Command{
			{
				Name:    "increment",
				Aliases: []string{"i"},
				Usage:   "increment the counter",
				Action: func(c *cli.Context) error {
					return runAction(c, func(params *Parameters) error {
						_, err := cl.NewClient(createURL(params.host)).IncrementCounterValue()
						return err
					})
				},
			},
			{
				Name:    "decrement",
				Aliases: []string{"d"},
				Usage:   "decrement the counter",
				Action: func(c *cli.Context) error {
					return runAction(c, func(params *Parameters) error {
						_, err := cl.NewClient(createURL(params.host)).DecrementCounterValue()
						return err
					})
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "get the counter value",
				Action: func(c *cli.Context) error {
					return runAction(c, func(params *Parameters) error {
						_, err := cl.NewClient(createURL(params.host)).GetCounterValue()
						return err
					})
				},
			},
		},
	}
}

//Parameters struct that conveys the main parameters of the app
type Parameters struct {
	host string
}

func runAction(c *cli.Context, f func(params *Parameters) error) error {
	log.SetFormatter(&log.JSONFormatter{})
	lvl, err := log.ParseLevel(c.String("log-level"))
	if err != nil {
		log.Errorf("unable to parse log level (%s): %v",
			c.String("log-level"), err)
		return err
	}
	log.SetLevel(lvl)

	host := c.String("host")

	return f(&Parameters{host: host})
}

func createURL(baseHost string) string {
	return fmt.Sprintf(protocolSchemaTemplate,
		path.Join(baseHost,
			"api",
			model.APIVersion,
			"counter"))
}
