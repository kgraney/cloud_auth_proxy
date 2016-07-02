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
	}

	app.Commands = []cli.Command{
		{
			Name:   "googleapis",
			Usage:  "Start a reverse proxy to googleapis.com",
			Action: googleapis.Main,
		},
	}

	app.Run(os.Args)
}
