package k8s

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

const containerName = "cmd"

// Start starts the command by attaching to the Job's Pod and streaming I/O.
func (k *k8s) Start() error {
	ctx := context.Background()

	// Wait for the Job's Pod to be created and Running
	podName, err := k.waitForPodRunning(ctx, 5*time.Minute)
	if err != nil {
		return err
	}

	req := k.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(k.jobNamespace).
		Name(podName).
		SubResource("attach").
		VersionedParams(&corev1.PodAttachOptions{
			Container: containerName,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k.restConfig, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("k8s: new attach executor: %w", err)
	}

	// Run attach stream in goroutine so Start() can return and Wait() can wait for Job completion
	go func() {
		_ = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:  k.stdin,
			Stdout: k.stdout,
			Stderr: k.stderr,
			Tty:    true,
		})
	}()

	return nil
}

// waitForPodRunning waits for the Job's Pod to be Running and returns its name.
func (k *k8s) waitForPodRunning(ctx context.Context, timeout time.Duration) (string, error) {
	selector := "job-name=" + k.jobName
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		pods, err := k.clientset.CoreV1().Pods(k.jobNamespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
		if err != nil {
			return "", fmt.Errorf("k8s: list pods: %w", err)
		}
		for _, p := range pods.Items {
			switch p.Status.Phase {
			case corev1.PodRunning:
				for _, cs := range p.Status.ContainerStatuses {
					if cs.Name == containerName && cs.Ready {
						return p.Name, nil
					}
				}
			case corev1.PodSucceeded, corev1.PodFailed:
				// Pod 已经完成（Succeeded 或 Failed），很可能是一次性 Job。
				// 这时再去 attach 也会很快 EOF；直接返回 Pod 名，
				// 让上层继续流程，避免在这里一直等 Running 而“卡住”。
				return p.Name, nil
			}
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(500 * time.Millisecond):
			continue
		}
	}
	return "", fmt.Errorf("k8s: timeout waiting for pod (job %s) to be running", k.jobName)
}
