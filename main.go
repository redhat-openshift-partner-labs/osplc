package main

import (
	"github.com/urfave/cli/v2"
	"os"

	"log"
)

func main() {
	app := &cli.App{
		Name:     "osplc",
		Usage:    "Manages ClusterDeployment and CronJob resources",
		Flags:    Flags,
		Commands: Commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
