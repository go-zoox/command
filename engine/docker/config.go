package docker

// Config is the configuration for a Docker engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string

	Image string
}
