package command

func (c *command) Cancel() error {
	return c.engine.Cancel()
}
