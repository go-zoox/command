# Agent Guide for go-zoox/command

This document provides implementation notes and checklists for maintainers and AI agents working on the command library (e.g. adding engines, changing config, or debugging).

## Engine Interface Contract

Every engine must implement `engine.Engine` in [engine/engine.go](engine/engine.go):

- **Start() error** – Start the command (attach streams, launch process/container/job). Must not block; streaming is typically done in goroutines.
- **Wait() error** – Block until the command finishes. Return `*errors.ExitError` for non-zero exit codes.
- **Cancel() error** – Terminate the command (kill process, remove container/job).
- **SetStdin / SetStdout / SetStderr** – Store the given `io.Reader`/`io.Writer`; use them when starting (e.g. in Start or attach).
- **Terminal() (terminal.Terminal, error)** – Return an interactive terminal implementing [terminal/terminal.go](terminal/terminal.go): `io.ReadWriteCloser`, `Resize(rows, cols)`, `ExitCode()`, `Wait()`.

## Terminal Interface

[terminal/terminal.go](terminal/terminal.go) requires:

- `Read`, `Write`, `Close` (io.ReadWriteCloser)
- `Resize(rows, cols int) error`
- `ExitCode() int`
- `Wait() error`

Engines that do not support TTY resize (e.g. k8s attach, wsl) may implement `Resize` as a no-op.

## Adding a New Engine

1. **Create `engine/<name>/`** with:
   - `<name>.go`: `Name` constant, `Config` (or use shared config), `New(cfg *Config) (engine.Engine, error)`, and a struct holding config + clients/handles.
   - `config.go`: Engine-specific config struct and defaults.
   - `create.go`: Allocate resources (e.g. create container/job); do not start yet unless the runtime requires it.
   - `start.go`: Attach stdin/stdout/stderr and start (or trigger start). Use goroutines for stream copying if needed.
   - `wait.go`: Wait for completion and map exit code to `*errors.ExitError`.
   - `cancel.go`: Clean up (kill process, remove container/job).
   - `io.go`: `SetStdin`, `SetStdout`, `SetStderr` (store references).
   - `terminal.go`: Return a type implementing `terminal.Terminal` (Resize may be no-op).

2. **Config**  
   In [config/config.go](config/config.go), add fields under a comment `// engine = <name>`.

3. **Registration**  
   In [init.go](init.go), call `engine.Register(<name>.Name, func(cfg *config.Config) (engine.Engine, error) { ... })` and map `config.Config` into the engine’s `Config`.

4. **Agent / command.go**  
   If the command can be executed via the agent, add the new engine’s config fields to the struct passed to `agent.New` in [command.go](command.go).

5. **Tests**
   - Unit tests in `engine/<name>/` (e.g. `config_test.go`, optional `create_test.go` with fakes).
   - In [command_test.go](command_test.go), add `TestEngine_<Name>` that uses `Engine: "<name>"` and runs a simple command (e.g. `echo hello`). If the engine depends on external services (Docker, K8s, Podman, WSL), **skip** when unavailable (e.g. check error message or env) so CI does not require that environment.

6. **Docs**  
   Update [README.md](README.md): add the engine to the “Execution Engines” list and add a short “<Engine> Configuration” example with the new config fields.

## Testing Conventions

- Tests that depend on Docker, Kubernetes, Podman, or WSL should **skip** when the backend is unavailable (e.g. `t.Skipf(...)` or `t.Skip(...)` on non-Windows for WSL).
- Prefer detecting “not available” from the error returned by `New` or `Run` (e.g. “docker”, “k8s”, “podman”, “wsl”) so the same test can run in CI without the service.
- Engine-specific config (e.g. Runtime, Kubeconfig) can be covered in `engine/<name>/config_test.go`.

## Engine-Specific Notes

- **K8s**: Uses a Job; create in `create()`, then in `Start()` wait for the Job’s Pod to be Running and attach via the attach API. Terminal() returns a limited terminal (Wait/ExitCode from Job status; Resize no-op).
- **Podman**: Reuses the Docker client with a different host (e.g. `unix:///run/podman/podman.sock`). Same flow as Docker: create container, attach, start, wait, remove.
- **WSL**: Only register/use on Windows (`runtime.GOOS == "windows"`). Build command as `wsl [-d Distro] -e <shell> -c "<command>"`; create `exec.Cmd` in Start() or Terminal(). Resize is no-op.
- **gVisor**: No new engine. Set `DockerRuntime: "runsc"` (and optionally `DockerRuntime: "kata"`) in config; [engine/docker/create.go](engine/docker/create.go) sets `hostCfg.Runtime` when `d.cfg.Runtime != ""`.

## File Change Checklist (new engine)

- [ ] `engine/<name>/*.go` (main, config, create, start, wait, cancel, io, terminal)
- [ ] `config/config.go` (new fields)
- [ ] `init.go` (Register + config mapping)
- [ ] `command.go` (agent config mapping if applicable)
- [ ] `engine/<name>/*_test.go`
- [ ] `command_test.go` (TestEngine_<Name> with skip when env missing)
- [ ] README.md (engines list + configuration example)
