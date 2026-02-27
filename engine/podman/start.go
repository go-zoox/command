package podman

import (
	"context"
	"io"
	"net"

	"github.com/docker/docker/api/types/container"
)

// Start starts the command.
func (p *podman) Start() error {
	stream, err := p.client.ContainerAttach(context.Background(), p.container.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}

	if err := applyStdin(stream.Conn, p.stdin); err != nil {
		return err
	}
	if err := applyStdout(stream.Conn, p.stdout); err != nil {
		return err
	}
	if err := applyStderr(stream.Conn, p.stderr); err != nil {
		return err
	}

	return p.client.ContainerStart(context.Background(), p.container.ID, container.StartOptions{})
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
