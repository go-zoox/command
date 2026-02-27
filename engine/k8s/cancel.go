package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Cancel deletes the Job (and its Pods via cascade).
func (k *k8s) Cancel() error {
	ctx := context.Background()
	propagation := metav1.DeletePropagationForeground
	return k.clientset.BatchV1().Jobs(k.jobNamespace).Delete(ctx, k.jobName, metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	})
}
