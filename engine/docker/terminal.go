package docker

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/go-zoox/command/errors"
	"github.com/go-zoox/command/terminal"
)

// Terminal returns a terminal.
func (d *docker) Terminal() (terminal.Terminal, error) {
	stream, err := d.client.ContainerAttach(context.Background(), d.container.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		// Logs:   true,
	})
	if err != nil {
		return nil, err
	}

	t := &Terminal{
		Ctx:         context.Background(),
		Client:      d.client,
		ContainerID: d.container.ID,
		Conn:        stream.Conn,
		ReadOnly:    d.cfg.ReadOnly,
	}

	err = d.client.ContainerStart(context.Background(), d.container.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Terminal is a terminal.
type Terminal struct {
	Ctx  context.Context
	Conn net.Conn
	//
	Client      *dockerClient.Client
	ContainerID string
	//
	ReadOnly bool
}

// Close closes the terminal.
func (t *Terminal) Close() error {
	t.Conn.Close()

	return t.Client.ContainerRemove(t.Ctx, t.ContainerID, types.ContainerRemoveOptions{
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

	return t.Conn.Write(p)
}

// Resize resizes the terminal.
func (t *Terminal) Resize(rows, cols int) error {
	inspect, err := t.Client.ContainerInspect(t.Ctx, t.ContainerID)
	if err != nil {
		return err
	}

	if inspect.State.Status != "running" {
		// return fmt.Errorf("container is not running")
		return nil
	}

	return t.Client.ContainerResize(t.Ctx, t.ContainerID, types.ResizeOptions{
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
			return fmt.Errorf("container exit error: %#v", err)
		}

	case result := <-resultC:
		if result.StatusCode != 0 {
			// return fmt.Errorf("container exited with non-zero status: %d", result.StatusCode)
			return &errors.ExitError{
				Code:    int(result.StatusCode),
				Message: fmt.Sprintf("container exited with non-zero status: %d", result.StatusCode),
			}
		}
	}

	return nil
}
