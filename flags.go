package main

import "github.com/urfave/cli/v2"

var _flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "cluster",
		Usage: "Name of the cluster",
	},
	&cli.StringFlag{
		Name:  "uptime",
		Usage: "Set/get the uptime for the cluster",
	},
	&cli.StringFlag{
		Name:  "timezone",
		Usage: "Set/get the timezone for the cluster",
	},
}
