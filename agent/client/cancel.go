package client

import (
	"time"

	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/core-utils/fmt"
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

	timer := time.NewTicker(30 * time.Second)
	select {
	case <-c.core.Context().Done():
		return c.core.Context().Err()
	case <-c.cancelEventDone:
		return nil
	case <-timer.C:
		return fmt.Errorf("timeout to wait start event")
	}
}
