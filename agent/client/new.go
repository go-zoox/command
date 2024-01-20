package client

import (
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

	<-c.newEventDone

	return nil
}
