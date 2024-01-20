package client

import (
	"fmt"
	"time"

	"github.com/go-zoox/command/agent/event"
	command "github.com/go-zoox/command/config"
	"github.com/go-zoox/logger"
)

func (c *client) New(command *command.Config) error {
	logger.Debugf("new event with command: %s", command.Command)

	err := c.sendEvent(&event.Event{
		Type:    event.New,
		Payload: command,
	})
	if err != nil {
		return err
	}

	timer := time.NewTicker(30 * time.Second)
	select {
	case <-c.core.Context().Done():
		return c.core.Context().Err()
	case <-c.newEventDone:
		return nil
	case <-timer.C:
		return fmt.Errorf("timeout to await new event")
	}
}
