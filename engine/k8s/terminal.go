package k8s

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	cmderrors "github.com/go-zoox/command/errors"
	"github.com/go-zoox/command/terminal"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// Terminal returns a terminal for the k8s job.
// It attaches to the Job's Pod using SPDY and exposes a ReadWriteCloser interface
// compatible with the command.Terminal usage in cmd/cmd.
func (k *k8s) Terminal() (terminal.Terminal, error) {
	ctx := context.Background()

	// Wait for Pod to be Running or already completed (for short-lived jobs).
	podName, err := k.waitForPodRunning(ctx, 5*time.Minute)
	if err != nil {
		return nil, err
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

	exec, err := remotecommand.NewSPDYExecutor(k.restConfig, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("k8s: new attach executor (terminal): %w", err)
	}

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	// Start attach stream in background.
	go func() {
		defer stdoutWriter.Close()
		defer stdinReader.Close()

		_ = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:  stdinReader,
			Stdout: stdoutWriter,
			Stderr: stdoutWriter, // merge stderr into stdout for simplicity
			Tty:    true,
		})
	}()

	return &Terminal{
		k8s:      k,
		stdin:    stdinWriter,
		stdout:   stdoutReader,
		ReadOnly: k.cfg.ReadOnly,
	}, nil
}

// Terminal is the terminal implementation for k8s.
type Terminal struct {
	k8s      *k8s
	ReadOnly bool

	stdin  io.WriteCloser
	stdout io.ReadCloser

	mu     sync.Mutex
	closed bool

	exitCode int
	waited   bool
}

// Read reads from the attached pod stdout/stderr stream.
func (t *Terminal) Read(p []byte) (n int, err error) {
	return t.stdout.Read(p)
}

// Write writes to the attached pod stdin stream (disabled when ReadOnly).
func (t *Terminal) Write(p []byte) (n int, err error) {
	if t.ReadOnly {
		return 0, nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return 0, io.EOF
	}
	return t.stdin.Write(p)
}

// Close closes stdin; the remote stream will terminate when the process exits.
func (t *Terminal) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.closed {
		return nil
	}
	t.closed = true
	_ = t.stdin.Close()
	return nil
}

// Resize is currently a no-op; implementing remote resize would require an
// additional exec/attach call with a TerminalSizeQueue.
func (t *Terminal) Resize(rows, cols int) error {
	return nil
}

// ExitCode returns the container exit code from the Job's Pod
// (after Wait has been called or job completed).
func (t *Terminal) ExitCode() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.waited {
		return t.exitCode
	}

	ctx := context.Background()
	pods, err := t.k8s.clientset.CoreV1().Pods(t.k8s.jobNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + t.k8s.jobName})
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

// Wait waits for the Job to complete and sets exit code (similar to k8s.Wait()).
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
		return fmt.Errorf("k8s: wait job (terminal): %w", err)
	}

	code := 0
	pods, _ := t.k8s.clientset.CoreV1().Pods(t.k8s.jobNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: "job-name=" + t.k8s.jobName})
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
		return &cmderrors.ExitError{
			Code:    code,
			Message: fmt.Sprintf("job exited with status %d", code),
		}
	}
	return nil
}
