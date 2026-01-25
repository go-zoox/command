package docker

// Config is the configuration for a Docker engine.
type Config struct {
	Command     string
	Environment map[string]string
	WorkDir     string
	User        string
	Shell       string
	// ReadOnly means none-interactive for terminal, which is used for show log, like top
	ReadOnly bool

	// engine = docker
	// Image is the name of the docker image
	Image string
	// Memory is the memory limit, unit: MB
	Memory int64
	// CPU is the CPU limit, unit: core
	CPU float64
	// Platform is the command platform, available: linux/amd64, linux/arm64
	Platform string
	// Network is the network name
	Network string
	// DisableNetwork disables network
	DisableNetwork bool
	//
	Privileged bool
	// DockerHost is the Docker host address
	DockerHost string
	// ImageRegistry is the Docker image registry address
	ImageRegistry string
	// ImageRegistryUsername is the Docker image registry username
	ImageRegistryUsername string
	// ImageRegistryPassword is the Docker image registry password
	ImageRegistryPassword string

	// Custom Command Runner ID
	ID string

	// AllowedSystemEnvKeys is the allowed system environment keys, which will be inherited to the command
	AllowedSystemEnvKeys []string

	// DataDirOuter is the outer data directory
	DataDirOuter string
	// DataDirInner is the inner data directory
	DataDirInner string

	// Sandbox enables strict security settings for untrusted code
	Sandbox bool
}
