package podman

// Config is the configuration for the podman engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string
	ReadOnly    bool

	Image          string
	Memory         int64
	CPU            float64
	Platform       string
	Network        string
	DisableNetwork bool
	Privileged     bool

	// PodmanHost is the Podman socket (Docker-compatible API). Default: unix:///run/podman/podman.sock
	PodmanHost string

	ID string

	AllowedSystemEnvKeys []string
}
