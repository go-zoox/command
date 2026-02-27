# Command - Better Cmd for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/command)](https://pkg.go.dev/github.com/go-zoox/command)
[![Build Status](https://github.com/go-zoox/command/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-zoox/command/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/command)](https://goreportcard.com/report/github.com/go-zoox/command)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/command/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/command?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/command.svg)](https://github.com/go-zoox/command/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/command.svg?label=Release)](https://github.com/go-zoox/command/tags)

A powerful and flexible command execution library for Go, supporting multiple execution engines (host, docker, dind, ssh, caas, k8s, podman, wsl) with built-in sandbox mode for secure execution of untrusted code.

## Installation

To install the package, run:
```bash
go get -u github.com/go-zoox/command
```

## Getting Started

### Basic Usage

```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-zoox/command"
)

func main() {
	// Simple command execution
	cfg := &command.Config{
		Command: "echo hello world",
	}

	cmd, err := command.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var buf strings.Builder
	cmd.SetStdout(&buf)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String()) // Output: hello world
}
```

### Using Docker Engine

```go
cfg := &command.Config{
	Command: "echo 'Running in Docker'",
	Engine:  "docker",
	Image:   "alpine:latest",
}

cmd, err := command.New(cfg)
if err != nil {
	log.Fatal(err)
}

if err := cmd.Run(); err != nil {
	log.Fatal(err)
}
```

### Sandbox Mode (Secure Execution)

Sandbox mode provides a secure environment for executing untrusted code with strict security settings:

```go
cfg := &command.Config{
	Command: "echo 'Running in sandbox'",
	Sandbox: true, // Enable sandbox mode
	// Automatically uses docker engine with strict security settings
}

cmd, err := command.New(cfg)
if err != nil {
	log.Fatal(err)
}

if err := cmd.Run(); err != nil {
	log.Fatal(err)
}
```

> 📖 **For detailed sandbox mode documentation, see [Sandbox Mode Guide](./docs/SANDBOX.md)**

## Features

### Execution Engines

- **host**: Execute commands directly on the host system (default)
- **docker**: Execute commands in Docker containers
- **dind**: Docker-in-Docker execution
- **ssh**: Execute commands on remote servers via SSH
- **caas**: Execute commands on Container-as-a-Service platforms
- **k8s**: Execute commands in Kubernetes (Job/Pod) within a cluster
- **podman**: Execute commands in Podman containers (Docker-compatible API)
- **wsl**: Execute commands via WSL on Windows (Windows only)

### Sandbox Mode

Sandbox mode is designed for executing untrusted code securely. When enabled, it:

- **Automatically uses Docker engine** - Forces execution in isolated containers
- **Applies strict security settings**:
  - Non-privileged mode (forced)
  - Read-only root filesystem
  - Capability restrictions (drops dangerous capabilities)
  - Network isolation (disabled by default)
  - Resource limits (512MB memory, 1 CPU core by default)
  - No new privileges allowed
  - Temporary directories mounted as tmpfs (noexec, nosuid)

#### Sandbox Mode Security Features

| Feature | Description |
|---------|-------------|
| **Non-privileged** | Container runs without root privileges |
| **Read-only rootfs** | Root filesystem is read-only, writable directories use tmpfs |
| **Capability dropping** | Drops all capabilities, only adds minimal necessary ones |
| **Network isolation** | Network disabled by default (can be enabled if needed) |
| **Resource limits** | Default: 512MB memory, 1 CPU core (configurable) |
| **Security options** | `no-new-privileges:true` prevents privilege escalation |

#### Sandbox Mode Example

```go
cfg := &command.Config{
	Command: "python3 -c 'print(\"Hello from sandbox\")'",
	Sandbox: true,
	// Optional: Override default resource limits
	Memory: 1024, // 1GB memory limit
	CPU:    2.0,  // 2 CPU cores
	// Optional: Enable network if needed
	DisableNetwork: false,
	Network:        "bridge",
}

cmd, err := command.New(cfg)
if err != nil {
	log.Fatal(err)
}

if err := cmd.Run(); err != nil {
	log.Fatal(err)
}
```

#### Sandbox Mode Requirements

- **Linux only**: Sandbox mode uses Linux-specific security features (seccomp, capabilities)
- **Docker required**: Sandbox mode requires Docker to be installed and running
- **Docker engine**: Sandbox mode automatically uses docker engine (cannot use other engines)
- **Optional gVisor**: Set `DockerRuntime: "runsc"` to use gVisor for stronger isolation (requires runsc installed on the host)

## Configuration

### Basic Configuration

```go
cfg := &command.Config{
	Command: "your command here",
	WorkDir: "/path/to/workdir",
	Shell:   "/bin/bash",
	User:    "username",
}
```

### Docker Configuration

```go
cfg := &command.Config{
	Command: "your command",
	Engine:  "docker",
	Image:   "alpine:latest",
	Memory:  512,  // MB
	CPU:     1.0,  // cores
	Platform: "linux/amd64",
	// Optional: use gVisor or Kata for stronger isolation
	DockerRuntime: "runsc",
}
```

### Kubernetes Configuration

```go
cfg := &command.Config{
	Command: "echo hello",
	Engine:  "k8s",
	K8sNamespace: "default",
	K8sImage:     "alpine:latest",
	K8sKubeconfig: "",  // optional; uses in-cluster or KUBECONFIG
}
```

### Podman Configuration

```go
cfg := &command.Config{
	Command:    "your command",
	Engine:     "podman",
	Image:      "alpine:latest",
	PodmanHost: "unix:///run/podman/podman.sock",  // optional
}
```

### WSL Configuration (Windows only)

```go
cfg := &command.Config{
	Command:   "echo hello",
	Engine:    "wsl",
	WSLDistro: "Ubuntu",  // optional distribution name
}
```

### Advanced Configuration

```go
cfg := &command.Config{
	Command: "your command",
	Engine:  "docker",
	Image:   "custom-image:tag",
	
	// Resource limits
	Memory: 1024, // MB
	CPU:    2.0,  // cores
	
	// Network configuration
	Network:        "custom-network",
	DisableNetwork: false,
	
	// Security
	Privileged: false,
	
	// Docker registry
	ImageRegistry:         "registry.example.com",
	ImageRegistryUsername: "username",
	ImageRegistryPassword: "password",
	
	// Data directories
	DataDirOuter: "/host/path",
	DataDirInner: "/container/path",

	// Optional: use gVisor (runsc) or Kata for stronger isolation
	DockerRuntime: "runsc",
}
```

## API Reference

### Creating a Command

```go
cmd, err := command.New(&command.Config{
	Command: "echo hello",
})
```

### Running Commands

```go
// Run and wait for completion
err := cmd.Run()

// Or start and wait separately
err := cmd.Start()
if err != nil {
	log.Fatal(err)
}
err = cmd.Wait()
```

### Capturing Output

```go
var stdout, stderr strings.Builder
cmd.SetStdout(&stdout)
cmd.SetStderr(&stderr)

err := cmd.Run()
if err != nil {
	log.Fatal(err)
}

fmt.Println("Output:", stdout.String())
fmt.Println("Error:", stderr.String())
```

### Getting Output Directly

```go
output, err := cmd.Output()
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(output))
```

## Examples

### Example 1: Basic Command Execution

```go
cfg := &command.Config{
	Command: "ls -la",
}

cmd, _ := command.New(cfg)
cmd.Run()
```

### Example 2: Docker Container Execution

```go
cfg := &command.Config{
	Command: "cat /etc/os-release",
	Engine:  "docker",
	Image:   "ubuntu:20.04",
}

cmd, _ := command.New(cfg)
cmd.Run()
```

### Example 3: Sandbox Mode for Untrusted Code

```go
cfg := &command.Config{
	Command: userProvidedCode, // Untrusted code
	Sandbox: true,              // Enable sandbox mode
}

cmd, err := command.New(cfg)
if err != nil {
	log.Fatal("Failed to create sandbox:", err)
}

if err := cmd.Run(); err != nil {
	log.Fatal("Command failed:", err)
}
```

### Example 4: Custom Resource Limits

```go
cfg := &command.Config{
	Command: "compute-intensive-task",
	Sandbox: true,
	Memory:  2048, // 2GB
	CPU:     4.0,  // 4 cores
}

cmd, _ := command.New(cfg)
cmd.Run()
```

## Security Considerations

### Sandbox Mode

- Sandbox mode is designed for executing **untrusted code**
- Always use sandbox mode when executing user-provided code
- Sandbox mode provides strong isolation but is not a replacement for proper security auditing
- Resource limits help prevent resource exhaustion attacks

### Best Practices

1. **Always use sandbox mode** for untrusted code execution
2. **Set appropriate resource limits** based on your use case
3. **Monitor container execution** and set timeouts
4. **Use read-only mounts** for data that shouldn't be modified
5. **Regularly update Docker images** to include security patches

## License

GoZoox is released under the [MIT License](./LICENSE).
