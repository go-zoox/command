package client

import (
	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/logger"
)

func (c *client) Cancel() error {
	logger.Debugf("cancel event")

	err := c.sendEvent(&event.Event{
		Type: event.Cancel,
	})

	if err != nil {
		return err
	}

	<-c.cancelEventDone

	return nil
}
