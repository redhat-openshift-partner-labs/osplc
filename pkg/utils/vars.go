package utils

import (
	"github.com/redhat-openshift-partner-labs/osplc/pkg/k8s"
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

var (
	DC dynamic.Interface
	KC *kubernetes.Clientset

	ClusterDeploymentGVR = schema.GroupVersionResource{
		Group:    "hive.openshift.io",
		Version:  "v1",
		Resource: "clusterdeployments",
	}
)

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
