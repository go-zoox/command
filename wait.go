package command

// Wait waits for the command to exit.
func (c *command) Wait() error {
	return c.engine.Wait()
}
