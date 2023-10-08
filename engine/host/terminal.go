package host

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/go-zoox/command/terminal"
)

func (h *host) Terminal() (terminal.Terminal, error) {
	terminal, err := pty.Start(h.cmd)
	if err != nil {
		return nil, err
	}

	return &Terminal{
		File: terminal,
		Cmd:  h.cmd,
	}, nil
}

type Terminal struct {
	*os.File
	Cmd *exec.Cmd
}

func (t *Terminal) Close() error {
	return t.File.Close()
}

func (t *Terminal) Read(p []byte) (n int, err error) {
	return t.File.Read(p)
}

func (t *Terminal) Write(p []byte) (n int, err error) {
	return t.File.Write(p)
}

func (rt *Terminal) Resize(rows, cols int) error {
	return pty.Setsize(rt.File, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}

func (rt *Terminal) ExitCode() int {
	return rt.Cmd.ProcessState.ExitCode()
}

func (rt *Terminal) Wait() error {
	return rt.Cmd.Wait()
}
