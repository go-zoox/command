package docker

import (
	"context"
	"io"
	"net"

	"github.com/docker/docker/api/types/container"
)

// Start starts the command.
func (d *docker) Start() error {
	stream, err := d.client.ContainerAttach(context.Background(), d.container.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		// Logs:   true,
	})
	if err != nil {
		return err
	}

	if err := applyStdin(stream.Conn, d.stdin); err != nil {
		return nil
	}

	if err := applyStdout(stream.Conn, d.stdout); err != nil {
		return nil
	}

	if err := applyStderr(stream.Conn, d.stderr); err != nil {
		return nil
	}

	err = d.client.ContainerStart(context.Background(), d.container.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	return nil
}

func applyStdin(conn net.Conn, stdin io.Reader) error {
	if stdin != nil {
		go io.Copy(conn, stdin)
	}

	return nil
}

func applyStdout(conn net.Conn, stdout io.Writer) error {
	if stdout != nil {
		go io.Copy(stdout, conn)
	}

	return nil
}

func applyStderr(conn net.Conn, stderr io.Writer) error {
	return nil
}
