package main

import (
	"github.com/codegangsta/cli"
	"github.com/kgraney/cloud_auth_proxy/googleapis"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "cloud_auth_proxy"
	app.Usage = "An authenticating, authorizing, and logging proxy for Public Cloud APIs."
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Usage: "The port to listen on",
		},
		cli.StringFlag{
			Name:  "certfile",
			Usage: "The SSL certificate file to use (e.g. server.pem)",
		},
		cli.StringFlag{
			Name:  "keyfile",
			Usage: "The SSL keyfile to use (e.g. server.key)",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "googleapis",
			Usage:  "Start a reverse proxy to googleapis.com",
			Action: googleapis.Main,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "credentials",
					Usage: "Google API credentials (e.g. \"Google Sandbox.json\")",
				},
			},
		},
	}

	app.Run(os.Args)
}
