package dind

// Start starts the command.
func (d *dind) Start() error {
	return d.client.Start()
}
