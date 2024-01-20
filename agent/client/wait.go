package client

import (
	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/command/errors"
	"github.com/go-zoox/logger"
)

func (c *client) Wait() error {
	logger.Debugf("wait event")

	err := c.sendEvent(&event.Event{
		Type: event.Wait,
	})
	if err != nil {
		return err
	}

	select {
	case <-c.core.Context().Done():
		return c.core.Context().Err()
	case <-c.waitEventDone:
		logger.Debugf("wait for exit code ...")
		code := <-c.exitcodeCh
		logger.Debugf("exit code is %d", code)

		if code == 0 {
			return nil
		}

		return &errors.ExitError{
			Code: code,
		}
	}
}
