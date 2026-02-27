package podman

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/go-zoox/command/errors"
	"github.com/go-zoox/command/terminal"
)

// Terminal returns a terminal.
func (p *podman) Terminal() (terminal.Terminal, error) {
	stream, err := p.client.ContainerAttach(context.Background(), p.container.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, err
	}

	t := &Terminal{
		Ctx:         context.Background(),
		Client:      p.client,
		ContainerID: p.container.ID,
		Conn:        stream.Conn,
		ReadOnly:    p.cfg.ReadOnly,
	}

	err = p.client.ContainerStart(context.Background(), p.container.ID, container.StartOptions{})
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Terminal is a terminal for podman.
type Terminal struct {
	Ctx  context.Context
	Conn net.Conn
	//
	Client      *dockerClient.Client
	ContainerID string
	//
	ReadOnly bool
	//
	sync.Mutex
}

// Close closes the terminal.
func (t *Terminal) Close() error {
	t.Conn.Close()
	return t.Client.ContainerRemove(t.Ctx, t.ContainerID, container.RemoveOptions{
		Force: true,
	})
}

// Read reads from the terminal.
func (t *Terminal) Read(p []byte) (n int, err error) {
	return t.Conn.Read(p)
}

// Write writes to the terminal.
func (t *Terminal) Write(p []byte) (n int, err error) {
	if t.ReadOnly {
		return 0, nil
	}
	t.Lock()
	defer t.Unlock()
	return t.Conn.Write(p)
}

// Resize resizes the terminal.
func (t *Terminal) Resize(rows, cols int) error {
	inspect, err := t.Client.ContainerInspect(t.Ctx, t.ContainerID)
	if err != nil {
		return err
	}
	if inspect.State.Status != "running" {
		return nil
	}
	return t.Client.ContainerResize(t.Ctx, t.ContainerID, container.ResizeOptions{
		Height: uint(rows),
		Width:  uint(cols),
	})
}

// ExitCode returns the exit code.
func (t *Terminal) ExitCode() int {
	inspect, err := t.Client.ContainerInspect(t.Ctx, t.ContainerID)
	if err != nil {
		return -1
	}
	return inspect.State.ExitCode
}

// Wait waits for the terminal to exit.
func (t *Terminal) Wait() error {
	resultC, errC := t.Client.ContainerWait(t.Ctx, t.ContainerID, container.WaitConditionNotRunning)
	select {
	case err := <-errC:
		if err != nil && err != io.EOF {
			return fmt.Errorf("podman: container wait: %w", err)
		}
	case result := <-resultC:
		if result.StatusCode != 0 {
			return &errors.ExitError{
				Code:    int(result.StatusCode),
				Message: fmt.Sprintf("container exited with non-zero status: %d", result.StatusCode),
			}
		}
	}
	return nil
}
