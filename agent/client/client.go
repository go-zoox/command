package client

import (
	"io"
	"os"

	"github.com/go-zoox/command"
	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket"
)

type Client interface {
	Connect() error
	Close() error
	//
	New(command *command.Config) error
	//
	Start() error
	Wait() error
	Cancel() error
	//
	SetStdin(stdin io.Reader) error
	SetStdout(stdout io.Writer) error
	SetStderr(stderr io.Writer) error
	//
	Run() error
}

type client struct {
	opt *Option
	//
	core websocket.Conn
	//
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	//
	exitcodeCh chan int
	//
	newEventDone    chan struct{}
	startEventDone  chan struct{}
	waitEventDone   chan struct{}
	cancelEventDone chan struct{}
}

type Option struct {
	Server string
}

func New(opts ...func(opt *Option)) Client {
	opt := &Option{
		Server: "ws://localhost:8080",
	}
	for _, o := range opts {
		o(opt)
	}

	return &client{
		opt: opt,
		//
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		//
		exitcodeCh: make(chan int),
		//
		newEventDone:    make(chan struct{}),
		startEventDone:  make(chan struct{}),
		waitEventDone:   make(chan struct{}),
		cancelEventDone: make(chan struct{}),
	}
}

func (c *client) sendEvent(evt *event.Event) error {
	s, err := evt.Encode()
	if err != nil {
		return err
	}

	logger.Debugf("send event to server: %s", s)
	if err := c.core.WriteTextMessage(s); err != nil {
		return err
	}

	return nil
}
