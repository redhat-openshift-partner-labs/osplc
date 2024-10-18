package utils

import (
	"github.com/redhat-openshift-partner-labs/osplc/pkg/cluster"
	"log"
)

var Commands = []*cli.Command{
	{
		Name:  "status",
		Usage: "Get current status of ClusterDeployments",
		Subcommands: []*cli.Command{
			{
				Name:  "all",
				Usage: "Get the status of all target resources",
				Action: func(c *cli.Context) error {
					log.Println(c)
					return statusAll()
				},
			},
		},
	},
	{
		Name:  "setup",
		Usage: "Setup the cluster deployments and cron jobs",
		Action: func(c *cli.Context) error {
			return runSetup()
		},
	},
	{
		Name: "set",
		Subcommands: []*cli.Command{
			{
				Name:  "uptime",
				Usage: "Set the uptime for the cluster",
				Action: func(c *cli.Context) error {
					return setUptime(c.String("uptime"))
				},
			},
			{
				Name:  "timezone",
				Usage: "Set the timezone for the cluster",
				Action: func(c *cli.Context) error {
					return setTimezone(c.String("timezone"))
				},
			},
		},
	},
	{
		Name:  "get",
		Usage: "Get the status of the cluster",
		Action: func(c *cli.Context) error {
			return getStatus("cluster", c.String("cluster"))
		},
	},
}

func runSetup() error {
	// Process ClusterDeployments using dynamic client
	return cluster.ProcessClusterDeployments(KC, DC)
}

func setUptime(uptime string) error {

	return nil
}

func setTimezone(timezone string) error {

	return nil
}

func statusAll() error {
	// Get ClusterDeployment details using dynamic client
	return cluster.GetClusterDeploymentStatus(DC, "")
}

func getStatus(kind, name string) error {

	if kind == "cluster" {
		log.Printf("Getting status for cluster %s\n", name)
		return cluster.GetClusterDeploymentStatus(DC, name)
	}

	return nil
}
