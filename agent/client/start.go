package client

import (
	"time"

	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/core-utils/fmt"
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

	timer := time.NewTicker(30 * time.Second)
	defer timer.Stop()

	select {
	case <-c.core.Context().Done():
		return c.core.Context().Err()
	case <-c.startEventDone:
		return nil
	case <-timer.C:
		return fmt.Errorf("timeout to wait start event")
	}
}
