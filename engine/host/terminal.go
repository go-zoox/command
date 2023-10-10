package host

import (
	"os"
	"os/exec"

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
}

// Close closes the terminal.
func (t *Terminal) Close() error {
	return t.File.Close()
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
