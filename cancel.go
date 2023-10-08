package command

// Cancel cancels the command.
func (c *command) Cancel() error {
	return c.engine.Cancel()
}
