package k8s

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/go-zoox/command/errors"
	"github.com/go-zoox/command/terminal"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Terminal returns a terminal for the k8s job. Resize is not supported; Wait/ExitCode use Job status.
func (k *k8s) Terminal() (terminal.Terminal, error) {
	return &Terminal{
		k8s:      k,
		ReadOnly: k.cfg.ReadOnly,
	}, nil
}

// Terminal is a terminal that supports Wait and ExitCode; Read/Write/Resize are limited.
type Terminal struct {
	k8s      *k8s
	ReadOnly bool
	mu       sync.Mutex
	closed   bool
	exitCode int
	waited   bool
}

// Read returns EOF (streaming is handled by the attach goroutine from Start).
func (t *Terminal) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

// Write is a no-op (stdin is already attached in Start).
func (t *Terminal) Write(p []byte) (n int, err error) {
	if t.ReadOnly {
		return 0, nil
	}
	return 0, nil
}

// Close is a no-op.
func (t *Terminal) Close() error {
	t.mu.Lock()
	t.closed = true
	t.mu.Unlock()
	return nil
}

// Resize is not supported for k8s attach; no-op.
func (t *Terminal) Resize(rows, cols int) error {
	return nil
}

// ExitCode returns the container exit code from the Job's Pod (after Wait has been called or job completed).
func (t *Terminal) ExitCode() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.waited {
		return t.exitCode
	}
	ctx := context.Background()
	pods, err := t.k8s.clientset.CoreV1().Pods(t.k8s.jobNamespace).List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + t.k8s.jobName})
	if err != nil || len(pods.Items) == 0 {
		return -1
	}
	for _, c := range pods.Items[0].Status.ContainerStatuses {
		if c.Name == containerName && c.State.Terminated != nil {
			return int(c.State.Terminated.ExitCode)
		}
	}
	return -1
}

// Wait waits for the Job to complete and sets exit code.
func (t *Terminal) Wait() error {
	ctx := context.Background()
	err := wait.PollUntilContextTimeout(ctx, 500*time.Millisecond, 1*time.Hour, true, func(ctx context.Context) (bool, error) {
		job, err := t.k8s.clientset.BatchV1().Jobs(t.k8s.jobNamespace).Get(ctx, t.k8s.jobName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, c := range job.Status.Conditions {
			if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("k8s: wait job: %w", err)
	}

	code := 0
	pods, _ := t.k8s.clientset.CoreV1().Pods(t.k8s.jobNamespace).List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + t.k8s.jobName})
	if len(pods.Items) > 0 {
		for _, c := range pods.Items[0].Status.ContainerStatuses {
			if c.Name == containerName && c.State.Terminated != nil {
				code = int(c.State.Terminated.ExitCode)
				break
			}
		}
	}

	t.mu.Lock()
	t.exitCode = code
	t.waited = true
	t.mu.Unlock()

	if code != 0 {
		return &errors.ExitError{
			Code:    code,
			Message: fmt.Sprintf("job exited with status %d", code),
		}
	}
	return nil
}
