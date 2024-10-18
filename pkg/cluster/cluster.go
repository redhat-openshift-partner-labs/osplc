package cluster

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"log"

	"github.com/redhat-openshift-partner-labs/osplc/pkg/cronjob"
	"github.com/redhat-openshift-partner-labs/osplc/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetClusterDeploymentStatus(dynamicClient dynamic.Interface, clustername string) error {
	// Define GroupVersionResource for ClusterDeployment
	clusterDeploymentGVR := schema.GroupVersionResource{
		Group:    "hive.openshift.io",
		Version:  "v1",
		Resource: "clusterdeployments",
	}

	log.Println("Getting all ClusterDeployments")
	// List all ClusterDeployments across all namespaces
	clusterDeployments, err := dynamicClient.Resource(clusterDeploymentGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// For each clusterDeployment in the list, extract the name, namespace, and annotations
	// and store them in a slice of unstructured.Unstructured objects
	var clusterDeploymentDetails []unstructured.Unstructured
	for _, cd := range clusterDeployments.Items {
		name := cd.GetName()
		namespace := cd.GetNamespace()
		annots := cd.GetAnnotations()

		clusterDeploymentDetails = append(clusterDeploymentDetails, unstructured.Unstructured{
			Object: map[string]interface{}{
				"name":        name,
				"namespace":   namespace,
				"annotations": annots,
			},
		})
	}

	// Print the details of each ClusterDeployment in a table format
	for _, cd := range clusterDeploymentDetails {
		name := cd.Object["name"]
		annots := cd.Object["annotations"].(map[string]string)
		state := cd.Object["status"]
		timezone, tzExists := annots["timezone"]
		if !tzExists {
			timezone = "none"
		}
		fmt.Printf("ClusterDeployment: %s, State: %s, Timezone: %s\n", name, state, timezone)
	}

	return nil
}

func ProcessClusterDeployments(clientset *kubernetes.Clientset, dynamicClient dynamic.Interface) error {
	// Define GroupVersionResource for ClusterDeployment
	clusterDeploymentGVR := schema.GroupVersionResource{
		Group:    "hive.openshift.io",
		Version:  "v1",
		Resource: "clusterdeployments",
	}

	log.Println("Getting all ClusterDeployments")
	// List all ClusterDeployments across all namespaces
	clusterDeployments, err := dynamicClient.Resource(clusterDeploymentGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, cd := range clusterDeployments.Items {
		name := cd.GetName()
		namespace := cd.GetNamespace()
		annots := cd.GetAnnotations()

		timezone, tzExists := annots["timezone"]
		uptime, uptimeExists := annots["uptime"]

		if !tzExists || !uptimeExists {
			continue // Skip if annotations are missing
		}

		// Validate timezone
		if !utils.IsValidTimezone(timezone) {
			fmt.Printf("Invalid timezone for ClusterDeployment %s/%s\n", namespace, name)
			continue
		}

		// Validate uptime
		if !utils.IsValidUptime(uptime) {
			fmt.Printf("Invalid uptime for ClusterDeployment %s/%s\n", namespace, name)
			continue
		}

		log.Printf("ClusterDeployment: %s, Timezone: %s, Uptime: %s\n", name, timezone, uptime)

		// Create CronJob
		if err := cronjob.CreateCronJob(clientset, name, timezone); err != nil {
			fmt.Printf("Failed to create CronJob for ClusterDeployment %s/%s: %v\n", namespace, name, err)
			continue
		}

		// Update hibernateAfter in spec
		if err := updateHibernateAfter(dynamicClient, &cd, uptime); err != nil {
			fmt.Printf("Failed to update hibernateAfter for ClusterDeployment %s/%s: %v\n", namespace, name, err)
			continue
		}
	}

	return nil
}

func updateHibernateAfter(dynamicClient dynamic.Interface, cd *unstructured.Unstructured, uptime string) error {
	// Extract spec
	spec, found, err := unstructured.NestedMap(cd.Object, "spec")
	if err != nil || !found {
		return fmt.Errorf("spec not found in ClusterDeployment %s/%s", cd.GetNamespace(), cd.GetName())
	}

	// Check if hibernateAfter exists and matches uptime
	hibernateAfter, found, err := unstructured.NestedString(spec, "hibernateAfter")
	if err != nil {
		return fmt.Errorf("error retrieving hibernateAfter: %v", err)
	}

	if found && hibernateAfter == uptime {
		// Already up-to-date
		return nil
	}

	// Update hibernateAfter
	if err := unstructured.SetNestedField(spec, uptime, "hibernateAfter"); err != nil {
		return fmt.Errorf("failed to set hibernateAfter: %v", err)
	}

	// Update the spec in the ClusterDeployment object
	if err := unstructured.SetNestedMap(cd.Object, spec, "spec"); err != nil {
		return fmt.Errorf("failed to update spec: %v", err)
	}

	// Define GroupVersionResource for ClusterDeployment
	clusterDeploymentGVR := schema.GroupVersionResource{
		Group:    "hive.openshift.io",
		Version:  "v1",
		Resource: "clusterdeployments",
	}

	// Update the ClusterDeployment resource
	_, err = dynamicClient.Resource(clusterDeploymentGVR).Namespace(cd.GetNamespace()).Update(context.TODO(), cd, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ClusterDeployment: %v", err)
	}

	fmt.Printf("Updated hibernateAfter for ClusterDeployment %s/%s\n", cd.GetNamespace(), cd.GetName())
	return nil
}
