package command

import (
	"time"

	"github.com/go-zoox/core-utils/fmt"
)

// Run runs the command.
func (c *command) Run() error {
	if err := c.Start(); err != nil {
		return err
	}

	if c.cfg.Timeout != 0 {
		done := make(chan error)
		go func() {
			done <- c.Wait()
		}()

		timeout := time.After(c.cfg.Timeout)
		select {
		case <-timeout:
			c.Cancel()
			return fmt.Errorf("timeout to run command (timeout: %s)", c.cfg.Timeout)
		case err := <-done:
			return err
		}
	}

	return c.Wait()
}
