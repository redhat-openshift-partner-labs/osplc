package utils

import "github.com/urfave/cli/v2"

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "cluster",
		Usage: "Set the cluster name to update or get the status of",
	},
	&cli.IntFlag{
		Name:  "uptime",
		Usage: "Set/get the uptime for the cluster",
	},
	&cli.StringFlag{
		Name:  "timezone",
		Usage: "Set/get the timezone for the cluster",
	},
}
