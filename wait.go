package command

func (c *command) Wait() error {
	return c.engine.Wait()
}
