package wsl

import (
	"io"
	"os/exec"
	"sync"

	"github.com/go-zoox/command/errors"
	"github.com/go-zoox/command/terminal"
)

// Terminal returns a terminal for the wsl process. Resize is a no-op on Windows.
func (w *wsl) Terminal() (terminal.Terminal, error) {
	w.cmd = exec.Command("wsl", w.args...)
	if err := applyEnv(w.cmd, w.cfg.Environment, w.cfg.AllowedSystemEnvKeys); err != nil {
		return nil, err
	}
	if w.cfg.WorkDir != "" {
		w.cmd.Dir = w.cfg.WorkDir
	}

	stdinPipe, _ := w.cmd.StdinPipe()
	stdoutPipe, _ := w.cmd.StdoutPipe()
	stderrPipe, _ := w.cmd.StderrPipe()

	if err := w.cmd.Start(); err != nil {
		return nil, err
	}

	return &Terminal{
		Cmd:      w.cmd,
		stdin:    stdinPipe,
		stdout:   stdoutPipe,
		stderr:   stderrPipe,
		ReadOnly: w.cfg.ReadOnly,
	}, nil
}

// Terminal implements terminal.Terminal for WSL (Resize is no-op).
type Terminal struct {
	Cmd      *exec.Cmd
	stdin    io.WriteCloser
	stdout   io.ReadCloser
	stderr   io.ReadCloser
	ReadOnly bool
	mu       sync.Mutex
}

// Read reads from stdout.
func (t *Terminal) Read(p []byte) (n int, err error) {
	return t.stdout.Read(p)
}

// Write writes to stdin (no-op if ReadOnly).
func (t *Terminal) Write(p []byte) (n int, err error) {
	if t.ReadOnly {
		return 0, nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.stdin.Write(p)
}

// Close closes stdin and waits for process; stdout/stderr are left to the caller.
func (t *Terminal) Close() error {
	_ = t.stdin.Close()
	return nil
}

// Resize is a no-op on WSL (not supported).
func (t *Terminal) Resize(rows, cols int) error {
	return nil
}

// ExitCode returns the process exit code.
func (t *Terminal) ExitCode() int {
	if t.Cmd.ProcessState == nil {
		return -1
	}
	return t.Cmd.ProcessState.ExitCode()
}

// Wait waits for the process and returns ExitError if non-zero.
func (t *Terminal) Wait() error {
	if err := t.Cmd.Wait(); err != nil {
		if v, ok := err.(*exec.ExitError); ok {
			return &errors.ExitError{
				Code:    v.ExitCode(),
				Message: v.Error(),
			}
		}
		return &errors.ExitError{Code: 1, Message: err.Error()}
	}
	return nil
}
