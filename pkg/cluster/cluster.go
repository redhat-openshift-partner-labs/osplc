package cluster

import (
	"bytes"
	"context"
	"fmt"
	. "github.com/redhat-openshift-partner-labs/osplc/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/redhat-openshift-partner-labs/osplc/pkg/cronjob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetClusterDeploymentStatus(dynamicClient dynamic.Interface, clustername string) error {
	// Buffer for template output
	var details bytes.Buffer

	log.Println("Getting all ClusterDeployments")
	// List all ClusterDeployments across all namespaces
	clusterDeployments, err := dynamicClient.Resource(ClusterDeploymentGVR).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	tmpl, err := template.New("clusters").Parse(ClusterDetailsTemplate)
	if err != nil {
		return err
	}

	// For each clusterDeployment in the list, extract the name, namespace, and annotations
	// and store them in a slice of unstructured.Unstructured objects
	for _, cd := range clusterDeployments.Items {
		name := cd.GetName()
		annots := cd.GetAnnotations()
		state := cd.Object["status"].(map[string]interface{})["powerState"]
		timezone, tzExists := annots["timezone"]
		if !tzExists {
			timezone = "Undefined"
		}
		uptime, uptimeExists := annots["uptime"]
		if !uptimeExists {
			uptime = "Undefined"
		}
		detail := map[string]interface{}{
			"Name":     name,
			"State":    state,
			"UpTime":   uptime,
			"TimeZone": timezone,
		}
		if err = tmpl.Execute(os.Stdout, detail); err != nil {
			return err
		}
		fmt.Println(details.String())
	}

	// Print the details of each ClusterDeployment in a table format
	/*for _, cd := range clusterDeploymentDetails {
		name := cd.Object["name"]
		annots := cd.Object["annotations"].(map[string]string)
		state := cd.Object["status"]
		timezone, tzExists := annots["timezone"]
		if !tzExists {
			timezone = "none"
		}
		fmt.Printf("ClusterDeployment: %s, State: %s, Timezone: %s\n", name, state, timezone)
	}*/

	return nil
}

func ProcessClusterDeployments(clientset *kubernetes.Clientset, dynamicClient dynamic.Interface) error {
	log.Println("Getting all ClusterDeployments")
	// List all ClusterDeployments across all namespaces
	clusterDeployments, err := dynamicClient.Resource(ClusterDeploymentGVR).List(context.TODO(), metav1.ListOptions{})
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
			fmt.Printf("Missing timezone and uptime annotations")
			continue // Skip if annotations are missing
		}

		// Validate timezone
		if !IsValidTimezone(timezone) {
			fmt.Printf("Invalid timezone for ClusterDeployment %s/%s\n", namespace, name)
			continue
		}

		// Validate uptime
		if !IsValidUptime(uptime) {
			fmt.Printf("Invalid uptime for ClusterDeployment %s/%s\n", namespace, name)
			continue
		}

		log.Printf("ClusterDeployment: %s, Timezone: %s, Uptime: %s\n", name, timezone, uptime)

		// Create CronJob
		if err = cronjob.CreateCronJob(clientset, name, timezone); err != nil {
			fmt.Printf("Failed to create CronJob for ClusterDeployment %s/%s: %v\n", namespace, name, err)
			continue
		}

		// Update hibernateAfter in spec
		if err = updateHibernateAfter(dynamicClient, name, uptime); err != nil {
			fmt.Printf("Failed to update hibernateAfter for ClusterDeployment %s/%s: %v\n", namespace, name, err)
			continue
		}
	}

	return nil
}

func StartCluster(dynamicClient dynamic.Interface, name string) error {
	// Get the ClusterDeployment object
	cd, err := dynamicClient.Resource(ClusterDeploymentGVR).Namespace(name).Get(context.TODO(), name, metav1.GetOptions{})
	//cd, err := dynamicClient.Resource(clusterDeploymentGVR).Get(context.TODO(), name, metav1.GetOptions{})
	log.Printf("%+v", cd)
	if err != nil {
		return fmt.Errorf("failed to get ClusterDeployment %s: %v", name, err)
	}

	// Extract spec
	spec, found, err := unstructured.NestedMap(cd.Object, "spec")
	if err != nil || !found {
		return fmt.Errorf("spec not found in ClusterDeployment %s", name)
	}

	// Check if hibernateAfter exists
	hibernateAfter, found, err := unstructured.NestedString(spec, "hibernateAfter")
	if err != nil {
		return fmt.Errorf("error retrieving hibernateAfter: %v", err)
	}

	if !found {
		return fmt.Errorf("hibernateAfter not found in ClusterDeployment %s", name)
	}

	// Validate hibernateAfter
	if !IsValidUptime(hibernateAfter) {
		return fmt.Errorf("invalid hibernateAfter value in ClusterDeployment %s: %s", name, hibernateAfter)
	}

	// Update powerState to "Running"
	if err := unstructured.SetNestedField(cd.Object, "Running", "spec", "powerState"); err != nil {
		return fmt.Errorf("failed to update powerState: %v", err)
	}

	// Update the ClusterDeployment resource
	_, err = dynamicClient.Resource(ClusterDeploymentGVR).Namespace(name).Update(context.TODO(), cd, metav1.UpdateOptions{})
	//_, err = dynamicClient.Resource(clusterDeploymentGVR).Update(context.TODO(), cd, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ClusterDeployment: %v", err)
	}

	fmt.Printf("Started ClusterDeployment %s\n", name)
	return nil
}

func SetUptime(dynamicClient dynamic.Interface, cluster string, uptime string) error {
	cd, err := dynamicClient.Resource(ClusterDeploymentGVR).Namespace(cluster).Get(context.TODO(), cluster, metav1.GetOptions{})
	if err != nil {
		return err
	}

	name := cd.GetName()
	namespace := cd.GetNamespace()

	// Validate uptime
	if !IsValidUptime(uptime) {
		return fmt.Errorf("invalid uptime for ClusterDeployment %s/%s", namespace, name)
	}

	// Update uptime annotation
	if err = updateUptimeAnnotation(dynamicClient, cluster, uptime); err != nil {
		return err
		//return fmt.Errorf("Failed to update uptime for ClusterDeployment %s/%s: %v\n", namespace, name, err)
	}

	// Update hibernateAfter in spec
	if err = updateHibernateAfter(dynamicClient, cluster, uptime); err != nil {
		return err
		//return fmt.Errorf("Failed to update hibernateAfter for ClusterDeployment %s/%s: %v\n", namespace, name, err)
	}

	return nil
}

func updateUptimeAnnotation(dynamicClient dynamic.Interface, cluster string, uptime string) error {
	cd, err := dynamicClient.Resource(ClusterDeploymentGVR).Namespace(cluster).Get(context.TODO(), cluster, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// Extract annotations
	annots := cd.GetAnnotations()

	// Check if uptime exists and matches the new value
	if annots["uptime"] == uptime {
		fmt.Println("uptime already up-to-date")
		return nil
	}

	// Update uptime annotation
	annots["uptime"] = uptime
	cd.SetAnnotations(annots)

	// Update the ClusterDeployment resource
	_, err = dynamicClient.Resource(ClusterDeploymentGVR).Namespace(cd.GetNamespace()).Update(context.TODO(), cd, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ClusterDeployment: %v", err)
	}

	fmt.Printf("Updated uptime annotation for ClusterDeployment %s/%s\n", cd.GetNamespace(), cd.GetName())
	return nil
}

func updateHibernateAfter(dynamicClient dynamic.Interface, cluster string, uptime string) error {
	cd, err := dynamicClient.Resource(ClusterDeploymentGVR).Namespace(cluster).Get(context.TODO(), cluster, metav1.GetOptions{})
	if err != nil {
		return err
	}

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
		fmt.Println("hibernateAfter already up-to-date")
		return nil
	}

	// Update hibernateAfter
	if err = unstructured.SetNestedField(spec, uptime, "hibernateAfter"); err != nil {
		return fmt.Errorf("failed to set hibernateAfter: %v", err)
	}

	// Update the spec in the ClusterDeployment object
	if err = unstructured.SetNestedMap(cd.Object, spec, "spec"); err != nil {
		return fmt.Errorf("failed to update spec: %v", err)
	}

	// Update the ClusterDeployment resource
	_, err = dynamicClient.Resource(ClusterDeploymentGVR).Namespace(cd.GetNamespace()).Update(context.TODO(), cd, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ClusterDeployment: %v", err)
	}

	fmt.Printf("Updated hibernateAfter for ClusterDeployment %s/%s\n", cd.GetNamespace(), cd.GetName())
	return nil
}

func getCurrentTimeInTimezone(timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Error loading timezone: %v\n", err)
		return ""
	}

	// Return current time in the specified timezone in HH:MM format with AM/PM
	return time.Now().In(loc).Format("03:04 PM")
}
