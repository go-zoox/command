package command

// NewOptions is the options for new
type NewOptions struct {
	Context     string
	Environment map[string]string
	Shell       string
}

// New creates a Command
func New(script string, options ...*NewOptions) *Command {
	var context string
	var environment map[string]string
	var shell = "/bin/sh"
	if len(options) > 0 && options[0] != nil {
		context = options[0].Context
		environment = options[0].Environment
	}

	return &Command{
		Script:      script,
		Context:     context,
		Environment: environment,
		Shell:       shell,
	}
}
