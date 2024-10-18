package main

import (
	"github.com/redhat-openshift-partner-labs/osplc/pkg/k8s"
	"github.com/redhat-openshift-partner-labs/osplc/pkg/utils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"os"

	"log"
)

var DC dynamic.Interface
var KC *kubernetes.Clientset

func init() {
	var err error

	// Setup clients for use in the commands
	DC, err = k8s.GetDynamicClient()
	if err != nil {
		_ = cli.Exit("Unable to create require dynamic client", 100)
	}

	KC, err = k8s.GetClientSet()
	if err != nil {
		_ = cli.Exit("Unable to create require kubernetes client", 101)
	}
}

func main() {
	app := &cli.App{
		Name:     "zone-handler",
		Usage:    "Manages ClusterDeployment and CronJob resources",
		Flags:    utils.Flags,
		Commands: utils.Commands,
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
