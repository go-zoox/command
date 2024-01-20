package client

func (c *client) Close() error {
	return c.core.Close()
}
