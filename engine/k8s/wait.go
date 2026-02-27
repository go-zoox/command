package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/go-zoox/command/errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Wait waits for the Job to complete and returns the container exit code as error if non-zero.
func (k *k8s) Wait() error {
	ctx := context.Background()

	var job *batchv1.Job
	err := wait.PollUntilContextTimeout(ctx, 500*time.Millisecond, 1*time.Hour, true, func(ctx context.Context) (bool, error) {
		var err error
		job, err = k.clientset.BatchV1().Jobs(k.jobNamespace).Get(ctx, k.jobName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, c := range job.Status.Conditions {
			if c.Type == batchv1.JobComplete && c.Status == corev1.ConditionTrue {
				return true, nil
			}
			if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("k8s: wait for job: %w", err)
	}

	// Get exit code from the Pod's container status
	exitCode := 0
	pods, err := k.clientset.CoreV1().Pods(k.jobNamespace).List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + k.jobName})
	if err == nil && len(pods.Items) > 0 {
		for _, c := range pods.Items[0].Status.ContainerStatuses {
			if c.Name == containerName && c.State.Terminated != nil {
				exitCode = int(c.State.Terminated.ExitCode)
				break
			}
		}
	}

	if exitCode != 0 {
		return &errors.ExitError{
			Code:    exitCode,
			Message: fmt.Sprintf("job %s exited with status %d", k.jobName, exitCode),
		}
	}
	return nil
}
