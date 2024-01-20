package client

func (c *client) Run() error {
	if err := c.Start(); err != nil {
		return err
	}

	return c.Wait()
}
