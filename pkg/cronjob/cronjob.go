package cronjob

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ptr "k8s.io/utils/pointer"
)

func CreateCronJob(clientset *kubernetes.Clientset, name, timezone string) error {
	cronJob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          "00 08 * * 1-5",
			TimeZone:          &timezone,
			ConcurrencyPolicy: batchv1.ReplaceConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					BackoffLimit: ptr.Int32(0),
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy:      "Never",
							ServiceAccountName: "openshift-partner-labs-osplc",
							Containers: []corev1.Container{
								{
									Name:    "handler",
									Image:   "quay.io/rhopl/osplc:v0.0.1",
									Command: []string{"/app/osplc", "start", "cluster", "--cluster", name},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clientset.BatchV1().CronJobs("default").Create(context.TODO(), cronJob, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Printf("CronJob %s already exists\n", name)
			return nil
		}
		return err
	}

	fmt.Printf("CronJob %s created successfully\n", name)
	return nil
}

func DeleteCronJob(clientset *kubernetes.Clientset, namespace string, name string) error {
	err := clientset.BatchV1().CronJobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			fmt.Printf("CronJob %s not found\n", name)
			return nil
		}
		return err
	}

	fmt.Printf("CronJob %s deleted successfully\n", name)
	return nil
}
