package command

import (
	"fmt"
	"time"
)

// Wait waits for the command to exit.
func (c *command) Wait() error {
	if c.cfg.Timeout != 0 {
		done := make(chan error)
		go func() {
			done <- c.engine.Wait()
		}()

		select {
		case <-c.cfg.Context.Done():
			return c.cfg.Context.Err()
		case <-time.After(c.cfg.Timeout):
			c.Cancel()
			return fmt.Errorf("timeout to run command (command: %s, timeout: %s)", c.cfg.Command, c.cfg.Timeout)
		case err := <-done:
			return err
		}
	}

	return c.engine.Wait()
}
