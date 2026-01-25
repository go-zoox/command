# Sandbox Mode Documentation

## Overview

Sandbox mode is a security feature designed for executing untrusted code in an isolated and secure environment. When enabled, it automatically configures Docker containers with strict security settings to prevent malicious code from affecting the host system.

## When to Use Sandbox Mode

Use sandbox mode when:
- Executing user-provided code
- Running untrusted scripts or applications
- Processing data from unknown sources
- Building multi-tenant applications
- Running code in shared environments

**Always use sandbox mode for untrusted code execution.**

## How It Works

When `Sandbox: true` is set in the configuration:

1. **Engine Selection**: Automatically switches to Docker engine (required)
2. **Security Configuration**: Applies strict security settings
3. **Resource Limits**: Sets default resource limits (configurable)
4. **Isolation**: Creates an isolated execution environment

## Security Features

### 1. Non-Privileged Mode

Containers run without root privileges, preventing access to sensitive system resources.

```go
cfg := &command.Config{
	Sandbox: true,
	// Privileged is automatically set to false
}
```

### 2. Read-Only Root Filesystem

The root filesystem is mounted as read-only. Writable directories (`/tmp`, `/var/tmp`) are mounted as tmpfs with security restrictions.

**Benefits:**
- Prevents modification of system files
- Prevents installation of persistent malware
- Limits attack surface

**Writable Directories:**
- `/tmp` - Temporary files (tmpfs, noexec, nosuid)
- `/var/tmp` - Temporary files (tmpfs, noexec, nosuid)

### 3. Capability Restrictions

Docker capabilities are restricted to the minimum necessary set:

**Dropped:** All capabilities (via `CapDrop: ["ALL"]`)

**Added (minimal set):**
- `CHOWN` - Change file ownership
- `DAC_OVERRIDE` - Bypass file read/write/execute permission checks
- `FOWNER` - Bypass permission checks on operations
- `FSETID` - Don't clear set-user-ID/set-group-ID bits
- `KILL` - Send signals
- `SETGID` - Manipulate process GIDs
- `SETUID` - Manipulate process UIDs
- `SETPCAP` - Set process capabilities
- `NET_BIND_SERVICE` - Bind to ports < 1024
- `NET_RAW` - Use RAW and PACKET sockets
- `SYS_CHROOT` - Use chroot()
- `MKNOD` - Create special files
- `AUDIT_WRITE` - Write to audit log
- `SETFCAP` - Set file capabilities

### 4. Network Isolation

Network access is disabled by default for maximum security.

```go
cfg := &command.Config{
	Sandbox: true,
	// DisableNetwork is automatically set to true
}
```

**To enable network access:**
```go
cfg := &command.Config{
	Sandbox:        true,
	DisableNetwork: false,
	Network:        "bridge", // or custom network
}
```

### 5. Resource Limits

Default resource limits prevent resource exhaustion attacks:

- **Memory**: 512MB (default)
- **CPU**: 1.0 core (default)

**Customize limits:**
```go
cfg := &command.Config{
	Sandbox: true,
	Memory:  2048, // 2GB
	CPU:     2.0,  // 2 cores
}
```

### 6. Security Options

Additional security options are applied:

- `no-new-privileges:true` - Prevents processes from gaining new privileges

## Configuration Examples

### Basic Sandbox

```go
cfg := &command.Config{
	Command: "python3 script.py",
	Sandbox: true,
}

cmd, err := command.New(cfg)
if err != nil {
	log.Fatal(err)
}

if err := cmd.Run(); err != nil {
	log.Fatal(err)
}
```

### Sandbox with Custom Resources

```go
cfg := &command.Config{
	Command: "data-processing-task",
	Sandbox: true,
	Memory:  2048, // 2GB
	CPU:     4.0,  // 4 cores
}

cmd, _ := command.New(cfg)
cmd.Run()
```

### Sandbox with Network Access

```go
cfg := &command.Config{
	Command:        "curl https://api.example.com",
	Sandbox:        true,
	DisableNetwork: false,
	Network:        "bridge",
}

cmd, _ := command.New(cfg)
cmd.Run()
```

### Sandbox with Custom Docker Image

```go
cfg := &command.Config{
	Command: "node app.js",
	Sandbox: true,
	Image:   "node:18-alpine",
}

cmd, _ := command.New(cfg)
cmd.Run()
```

### Sandbox with Data Directory

```go
cfg := &command.Config{
	Command:     "process-data.sh",
	Sandbox:     true,
	DataDirOuter: "/host/data",
	DataDirInner: "/container/data",
}

cmd, _ := command.New(cfg)
cmd.Run()
```

## Default Behavior

When sandbox mode is enabled, the following defaults are applied:

| Setting | Default Value | Overridable |
|---------|--------------|-------------|
| Engine | `docker` | No (required) |
| Privileged | `false` | No (forced) |
| DisableNetwork | `true` | Yes |
| Memory | `512` MB | Yes |
| CPU | `1.0` core | Yes |
| ReadonlyRootfs | `true` | No (forced) |
| SecurityOpt | `no-new-privileges:true` | No (forced) |

## Limitations

### 1. Linux Only

Sandbox mode uses Linux-specific security features:
- Seccomp profiles
- Linux capabilities
- Namespaces

**macOS/Windows**: Sandbox mode will not work on macOS or Windows (Docker Desktop uses Linux VMs, but sandbox features may be limited).

### 2. Docker Required

Sandbox mode requires:
- Docker installed and running
- Docker daemon accessible
- Appropriate permissions to create containers

### 3. Read-Only Filesystem

Some applications may fail if they expect to write to the root filesystem. Use:
- `/tmp` or `/var/tmp` for temporary files
- `DataDirOuter`/`DataDirInner` for persistent data

### 4. Network Isolation

Network access is disabled by default. Enable it explicitly if needed:
```go
cfg.DisableNetwork = false
cfg.Network = "bridge"
```

## Security Best Practices

### 1. Always Use Sandbox Mode for Untrusted Code

```go
// ❌ Bad: Running untrusted code without sandbox
cfg := &command.Config{
	Command: userCode,
}

// ✅ Good: Using sandbox mode
cfg := &command.Config{
	Command: userCode,
	Sandbox: true,
}
```

### 2. Set Appropriate Resource Limits

```go
cfg := &command.Config{
	Command: userCode,
	Sandbox: true,
	Memory:  1024, // Based on expected usage
	CPU:     2.0,  // Based on expected load
}
```

### 3. Use Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

cfg := &command.Config{
	Command: userCode,
	Sandbox: true,
	Context: ctx,
	Timeout: 30 * time.Second,
}
```

### 4. Monitor Execution

```go
var stdout, stderr strings.Builder
cmd.SetStdout(&stdout)
cmd.SetStderr(&stderr)

err := cmd.Run()
if err != nil {
	log.Printf("Command failed: %v", err)
	log.Printf("Stdout: %s", stdout.String())
	log.Printf("Stderr: %s", stderr.String())
}
```

### 5. Keep Docker Images Updated

Regularly update Docker images to include security patches:

```go
cfg := &command.Config{
	Sandbox: true,
	Image:   "python:3.11", // Use latest stable version
}
```

## Troubleshooting

### Error: "sandbox mode requires docker engine"

**Cause**: Sandbox mode only works with Docker engine.

**Solution**: Don't specify a different engine when using sandbox mode:
```go
// ❌ Wrong
cfg := &command.Config{
	Sandbox: true,
	Engine:  "host", // Error!
}

// ✅ Correct
cfg := &command.Config{
	Sandbox: true,
	// Engine will be automatically set to "docker"
}
```

### Error: "Cannot connect to Docker daemon"

**Cause**: Docker is not running or not accessible.

**Solution**: 
1. Start Docker daemon: `sudo systemctl start docker`
2. Check Docker is running: `docker ps`
3. Verify permissions: Add user to docker group

### Application Fails: "Read-only file system"

**Cause**: Application tries to write to root filesystem.

**Solution**: Use writable directories:
```go
cfg := &command.Config{
	Command: "your-command",
	Sandbox: true,
	// Use /tmp or /var/tmp for temporary files
	// Use DataDirOuter/DataDirInner for persistent data
}
```

### Network Requests Fail

**Cause**: Network is disabled by default in sandbox mode.

**Solution**: Enable network explicitly:
```go
cfg := &command.Config{
	Sandbox:        true,
	DisableNetwork: false,
	Network:        "bridge",
}
```

## Performance Considerations

Sandbox mode adds some overhead:
- Container creation time
- Resource isolation overhead
- Security checks overhead

For trusted code, use host engine for better performance:
```go
// Trusted code - use host engine
cfg := &command.Config{
	Command: trustedCode,
	Engine:  "host",
}

// Untrusted code - use sandbox mode
cfg := &command.Config{
	Command: untrustedCode,
	Sandbox: true,
}
```

## Comparison: Sandbox vs Non-Sandbox

| Feature | Sandbox Mode | Non-Sandbox Mode |
|---------|-------------|------------------|
| Isolation | Strong (container) | None (host) |
| Security | High | Low |
| Performance | Lower (container overhead) | Higher (direct execution) |
| Resource Limits | Yes | No (unless configured) |
| Network Access | Disabled by default | Enabled |
| Root Filesystem | Read-only | Writable |
| Capabilities | Restricted | Full |
| Use Case | Untrusted code | Trusted code |

## Additional Resources

- [Docker Security Best Practices](https://docs.docker.com/engine/security/)
- [Linux Capabilities](https://man7.org/linux/man-pages/man7/capabilities.7.html)
- [Seccomp](https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt)
