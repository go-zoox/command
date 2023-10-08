package docker

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/go-zoox/command/terminal"
)

func (d *docker) Terminal() (terminal.Terminal, error) {
	stream, err := d.client.ContainerAttach(d.ctx, d.container.ID, types.ContainerAttachOptions{
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
		Ctx:         d.ctx,
		Client:      d.client,
		ContainerID: d.container.ID,
		Conn:        stream.Conn,
	}

	err = d.client.ContainerStart(d.ctx, d.container.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	return t, nil
}

type Terminal struct {
	Ctx  context.Context
	Conn net.Conn
	//
	Client      *dockerClient.Client
	ContainerID string
}

func (t *Terminal) Close() error {
	t.Conn.Close()

	return t.Client.ContainerRemove(t.Ctx, t.ContainerID, types.ContainerRemoveOptions{
		Force: true,
	})
}

func (t *Terminal) Read(p []byte) (n int, err error) {
	return t.Conn.Read(p)
}

func (t *Terminal) Write(p []byte) (n int, err error) {
	return t.Conn.Write(p)
}

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

func (t *Terminal) ExitCode() int {
	inspect, err := t.Client.ContainerInspect(t.Ctx, t.ContainerID)
	if err != nil {
		return -1
	}

	return inspect.State.ExitCode
}

func (rt *Terminal) Wait() error {
	resultC, errC := rt.Client.ContainerWait(rt.Ctx, rt.ContainerID, container.WaitConditionNotRunning)
	select {
	case err := <-errC:
		if err != nil && err != io.EOF {
			return fmt.Errorf("container exit error: %#v", err)
		}

	case result := <-resultC:
		if result.StatusCode != 0 {
			// rt.exitCode = int(result.StatusCode)
			return fmt.Errorf("container exited with non-zero status: %d", result.StatusCode)
		}
	}

	return nil
}
