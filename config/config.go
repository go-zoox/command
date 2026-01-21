package config

import (
	"context"
	"time"
)

// Config is the command config
type Config struct {
	Context context.Context

	// Timeout is the command timeout
	Timeout time.Duration

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
	IsHistoryDisabled           bool
	IsInheritEnvironmentEnabled bool
	//
	AllowedSystemEnvKeys []string

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
	// ImageRegistry is the Docker image registry address
	ImageRegistry string
	// ImageRegistryUsername is the Docker image registry username
	ImageRegistryUsername string
	// ImageRegistryPassword is the Docker image registry password
	ImageRegistryPassword string

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

	// Agent is the command runner agent address
	Agent string

	// DataDirOuter is the outer data directory
	DataDirOuter string
	// DataDirInner is the inner data directory
	DataDirInner string
}
