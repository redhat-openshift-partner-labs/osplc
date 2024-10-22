package main

import (
	"github.com/redhat-openshift-partner-labs/osplc/pkg/cluster"
	"github.com/redhat-openshift-partner-labs/osplc/pkg/cronjob"
	. "github.com/redhat-openshift-partner-labs/osplc/pkg/utils"
	"github.com/urfave/cli/v2"
	"log"
)

var _commands = []*cli.Command{
	{
		Name:  "status",
		Usage: "Get current status of ClusterDeployments",
		Subcommands: []*cli.Command{
			{
				Name:  "all",
				Usage: "Get the status of all target resources",
				Action: func(c *cli.Context) error {
					return statusAll()
				},
			},
		},
	},
	{
		Name:  "create",
		Usage: "Create resources",
		Subcommands: []*cli.Command{
			{
				Name:  "cronjob",
				Flags: _flags,
				Action: func(c *cli.Context) error {
					return cronjob.CreateCronJob(KC, c.String("cluster"), c.String("timezone"))
				},
			},
		},
	},
	{
		Name:  "destroy",
		Usage: "Destroy resources",
		Subcommands: []*cli.Command{
			{
				Name: "cronjob",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "Name of the cronjob",
					},
					&cli.StringFlag{
						Name:  "namespace",
						Usage: "Namespace of the cronjob",
					},
				},
				Action: func(c *cli.Context) error {
					return cronjob.DeleteCronJob(KC, c.String("namespace"), c.String("name"))
				},
			},
		},
	},
	{
		Name:  "start",
		Usage: "Start a resource",
		Subcommands: []*cli.Command{
			{
				Name:  "cluster",
				Usage: "Start a cluster",
				Flags: _flags,
				Action: func(c *cli.Context) error {
					return startCluster(c.String("cluster"))
				},
			},
		},
	},
	{
		Name: "set",
		Subcommands: []*cli.Command{
			{
				Name:  "uptime",
				Usage: "Set the uptime for the cluster; format is like 8h30m default is 8h",
				Flags: _flags,
				Action: func(c *cli.Context) error {
					return setUptime(c.String("cluster"), c.String("uptime"))
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
		Subcommands: []*cli.Command{
			{
				Name:  "cluster",
				Usage: "Get the status of a cluster",
				Flags: _flags,
				Action: func(c *cli.Context) error {
					return getStatus("cluster", c.String("cluster"))
				},
			},
		},
	},
}

func runSetup() error {
	// Process ClusterDeployments using dynamic client
	return cluster.ProcessClusterDeployments(KC, DC)
}

func setUptime(name, uptime string) error {
	return cluster.SetUptime(DC, name, uptime)
}

func setTimezone(timezone string) error {
	return nil
}

func startCluster(name string) error {
	// Start a cluster using dynamic client
	return cluster.StartCluster(DC, name)
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
