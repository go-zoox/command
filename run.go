package command

func (c *command) Run() error {
	if err := c.Start(); err != nil {
		return err
	}

	return c.Wait()
}
