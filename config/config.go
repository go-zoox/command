package config

import "context"

// Config is the command config
type Config struct {
	Context context.Context

	// Engine is the command engine, available: host, docker
	Engine string

	// engine common
	Command     string
	WorkDir     string
	Environment map[string]string
	User        string
	Shell       string
	// ReadOnly means none-interactive for terminal, which is used for show log, like top
	ReadOnly bool

	// engine = host
	IsHistoryDisabled bool

	// engine = docker
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
	// Privileged enables privileged mode
	Privileged bool
	// DockerHost is the Docker host
	DockerHost string

	// engine = caas
	// Server is the command server address
	Server string
	// ClientID is the client ID for server auth
	ClientID string
	// ClientSecret is the client secret for server auth
	ClientSecret string

	// engine = ssh
	SSHHost                          string
	SSHPort                          int
	SSHUser                          string
	SSHPass                          string
	SSHPrivateKey                    string
	SSHPrivateKeySecret              string
	SSHIsIgnoreStrictHostKeyChecking bool
	SSHKnowHostsFilePath             string

	// Custom Command Runner ID
	ID string
}
