package client

import (
	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/logger"
)

func (c *client) Start() error {
	logger.Debugf("start event")

	err := c.sendEvent(&event.Event{
		Type: event.Start,
	})
	if err != nil {
		return err
	}

	<-c.startEventDone

	return nil
}
