package host

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/creack/pty"
	"github.com/go-zoox/command/terminal"
)

// Name is the name of the engine.
func (h *host) Terminal() (terminal.Terminal, error) {
	terminal, err := pty.Start(h.cmd)
	if err != nil {
		return nil, err
	}

	return &Terminal{
		File:     terminal,
		Cmd:      h.cmd,
		ReadOnly: h.cfg.ReadOnly,
	}, nil
}

// Terminal is the terminal implementation.
type Terminal struct {
	*os.File
	Cmd      *exec.Cmd
	ReadOnly bool
	//
	sync.Mutex
	closeOnce sync.Once
	closeErr  error
}

// Close closes the PTY and tears down the PTY session. For /bin/sh -c "cmd", only
// killing the shell PID leaves children alive; signal the whole process group
// (negative PID == PGID, which matches the session leader from pty.Start).
func (t *Terminal) Close() error {
	t.closeOnce.Do(func() {
		if t.File != nil {
			if err := t.File.Close(); err != nil {
				t.closeErr = err
			}
		}
		if t.Cmd != nil && t.Cmd.Process != nil {
			pid := t.Cmd.Process.Pid
			// Kill the process group so e.g. cursor-agent outlives a dead shell.
			if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
				errno, _ := err.(syscall.Errno)
				if errno != syscall.ESRCH && errno != syscall.EINVAL {
					if t.closeErr == nil {
						t.closeErr = err
					}
				}
				if err := t.Cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
					if t.closeErr == nil {
						t.closeErr = err
					}
				}
			}
		}
	})
	return t.closeErr
}

// Read reads from the terminal.
func (t *Terminal) Read(p []byte) (n int, err error) {
	return t.File.Read(p)
}

// Write writes to the terminal.
func (t *Terminal) Write(p []byte) (n int, err error) {
	if t.ReadOnly {
		return 0, nil
	}
	t.Lock()
	defer t.Unlock()

	return t.File.Write(p)
}

// Resize resizes the terminal.
func (t *Terminal) Resize(rows, cols int) error {
	return pty.Setsize(t.File, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}

// ExitCode returns the exit code.
func (t *Terminal) ExitCode() int {
	return t.Cmd.ProcessState.ExitCode()
}

// Wait waits for the terminal to exit.
func (t *Terminal) Wait() error {
	return t.Cmd.Wait()
}
