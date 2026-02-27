package command

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	cmderrors "github.com/go-zoox/command/errors"
)

func TestNew(t *testing.T) {
	cfg := &Config{
		Command: "echo hello world",
	}

	cmd, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	buf := &strings.Builder{}
	cmd.SetStdout(buf)

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to start command: %v", err)
	}

	if v := buf.String(); v != "hello world\n" {
		t.Fatalf("expected %q, got %q", "hello world\n", v)
	}
}

func TestNew_DefaultEngine(t *testing.T) {
	cfg := &Config{
		Command: "echo ok",
		Engine:  "",
	}
	cmd, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	if cfg.Engine != "host" {
		t.Errorf("expected default engine host, got %q", cfg.Engine)
	}
	_ = cmd
}

func TestNew_UnsupportedEngine(t *testing.T) {
	cfg := &Config{
		Command: "echo ok",
		Engine:  "nonexistent-engine",
	}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for unsupported engine")
	}
	if !strings.Contains(err.Error(), "unsupported") && !strings.Contains(err.Error(), "nonexistent-engine") {
		t.Errorf("expected error to mention unsupported or engine name, got %q", err.Error())
	}
}

func TestNew_IDGenerated(t *testing.T) {
	cfg := &Config{
		Command: "echo ok",
		ID:     "",
	}
	_, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	if cfg.ID == "" {
		t.Error("expected ID to be generated when empty")
	}
	if !strings.HasPrefix(cfg.ID, "go-zoox_command_") {
		t.Errorf("expected ID prefix go-zoox_command_, got %q", cfg.ID)
	}
}

func TestNew_ShellDefault(t *testing.T) {
	cfg := &Config{
		Command: "echo ok",
		Shell:   "",
	}
	_, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	if cfg.Shell != "/bin/sh" {
		t.Errorf("expected default shell /bin/sh, got %q", cfg.Shell)
	}
}

func TestOutput(t *testing.T) {
	cfg := &Config{
		Command: "echo hello from output",
	}
	cmd, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("Output() failed: %v", err)
	}
	if v := strings.TrimSpace(string(out)); v != "hello from output" {
		t.Errorf("expected output %q, got %q", "hello from output", v)
	}
}

func TestStartThenWait(t *testing.T) {
	cfg := &Config{
		Command: "echo started then wait",
	}
	cmd, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	buf := &strings.Builder{}
	cmd.SetStdout(buf)
	if err := cmd.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}
	if v := buf.String(); !strings.Contains(v, "started then wait") {
		t.Errorf("expected stdout to contain 'started then wait', got %q", v)
	}
}

func TestSetStdin(t *testing.T) {
	cfg := &Config{
		Command: "cat",
	}
	cmd, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	cmd.SetStdin(strings.NewReader("stdin content\n"))
	buf := &strings.Builder{}
	cmd.SetStdout(buf)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Run() failed: %v", err)
	}
	if v := buf.String(); v != "stdin content\n" {
		t.Errorf("expected stdout from stdin, got %q", v)
	}
}

func TestSetStderr(t *testing.T) {
	cfg := &Config{
		Command: "echo to stderr 1>&2",
	}
	cmd, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	stderrBuf := &strings.Builder{}
	cmd.SetStderr(stderrBuf)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Run() failed: %v", err)
	}
	if v := stderrBuf.String(); !strings.Contains(v, "to stderr") {
		t.Errorf("expected stderr to contain 'to stderr', got %q", v)
	}
}

func TestRun_ExitNonZero(t *testing.T) {
	cfg := &Config{
		Command: "exit 42",
	}
	cmd, err := New(cfg)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}
	err = cmd.Run()
	if err == nil {
		t.Fatal("expected error for exit 42")
	}
	var exitErr *cmderrors.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected *errors.ExitError, got %T: %v", err, err)
	}
	if exitErr.ExitCode() != 42 {
		t.Errorf("expected exit code 42, got %d", exitErr.ExitCode())
	}
}

func TestSandboxMode_AutoSwitchToDocker(t *testing.T) {
	cfg := &Config{
		Command: "echo test",
		Sandbox: true,
		Engine:  "", // Empty engine should auto-switch to docker
	}

	// This will fail if docker is not available, but we can test the config logic
	_, err := New(cfg)
	if err != nil {
		// If docker is not available, that's expected
		if !strings.Contains(err.Error(), "docker") && !strings.Contains(err.Error(), "Docker") {
			t.Logf("Docker not available, skipping test: %v", err)
			return
		}
	}

	// Verify engine was switched to docker
	if cfg.Engine != "docker" {
		t.Errorf("expected engine to be 'docker', got %q", cfg.Engine)
	}
}

func TestSandboxMode_ForceDockerEngine(t *testing.T) {
	cfg := &Config{
		Command: "echo test",
		Sandbox: true,
		Engine:  "docker", // Explicitly set to docker
	}

	_, err := New(cfg)
	if err != nil {
		// If docker is not available, that's expected
		if !strings.Contains(err.Error(), "docker") && !strings.Contains(err.Error(), "Docker") {
			t.Logf("Docker not available, skipping test: %v", err)
			return
		}
	}

	// Verify engine remains docker
	if cfg.Engine != "docker" {
		t.Errorf("expected engine to be 'docker', got %q", cfg.Engine)
	}
}

func TestSandboxMode_RejectNonDockerEngine(t *testing.T) {
	cfg := &Config{
		Command: "echo test",
		Sandbox: true,
		Engine:  "host", // Try to use host engine with sandbox
	}

	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error when using non-docker engine with sandbox mode")
	}

	expectedError := "sandbox mode requires docker engine"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error to contain %q, got %q", expectedError, err.Error())
	}
}

func TestSandboxMode_DefaultSecuritySettings(t *testing.T) {
	cfg := &Config{
		Command:      "echo test",
		Sandbox:      true,
		Engine:       "docker",
		Privileged:   true,  // Should be forced to false
		Memory:       0,     // Should get default 512MB
		CPU:          0,     // Should get default 1.0
		DisableNetwork: false, // Should be set to true
		Network:      "",    // Empty network
	}

	// Create command to trigger sandbox logic
	_, err := New(cfg)
	if err != nil {
		// If docker is not available, skip the test
		if strings.Contains(err.Error(), "docker") || strings.Contains(err.Error(), "Docker") {
			t.Logf("Docker not available, skipping test: %v", err)
			return
		}
		// Other errors are unexpected
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify security settings were applied
	if cfg.Privileged {
		t.Error("expected Privileged to be false in sandbox mode")
	}

	if cfg.Memory != 512 {
		t.Errorf("expected default Memory to be 512MB, got %d", cfg.Memory)
	}

	if cfg.CPU != 1.0 {
		t.Errorf("expected default CPU to be 1.0, got %f", cfg.CPU)
	}

	if !cfg.DisableNetwork {
		t.Error("expected DisableNetwork to be true in sandbox mode when Network is empty")
	}
}

func TestSandboxMode_PreserveExplicitSettings(t *testing.T) {
	cfg := &Config{
		Command:      "echo test",
		Sandbox:      true,
		Engine:       "docker",
		Memory:       1024,  // Explicitly set
		CPU:          2.0,    // Explicitly set
		DisableNetwork: false, // Explicitly set
		Network:      "custom-network", // Explicitly set
	}

	_, err := New(cfg)
	if err != nil {
		// If docker is not available, skip the test
		if strings.Contains(err.Error(), "docker") || strings.Contains(err.Error(), "Docker") {
			t.Logf("Docker not available, skipping test: %v", err)
			return
		}
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify explicit settings were preserved
	if cfg.Memory != 1024 {
		t.Errorf("expected Memory to be 1024MB, got %d", cfg.Memory)
	}

	if cfg.CPU != 2.0 {
		t.Errorf("expected CPU to be 2.0, got %f", cfg.CPU)
	}

	// When Network is explicitly set, DisableNetwork should not be forced
	if cfg.DisableNetwork && cfg.Network != "" {
		t.Error("DisableNetwork should not be forced when Network is explicitly set")
	}
}

func TestSandboxMode_ForceNonPrivileged(t *testing.T) {
	cfg := &Config{
		Command:    "echo test",
		Sandbox:    true,
		Engine:     "docker",
		Privileged: true, // Try to enable privileged mode
	}

	_, err := New(cfg)
	if err != nil {
		// If docker is not available, skip the test
		if strings.Contains(err.Error(), "docker") || strings.Contains(err.Error(), "Docker") {
			t.Logf("Docker not available, skipping test: %v", err)
			return
		}
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify privileged mode was forced to false
	if cfg.Privileged {
		t.Error("expected Privileged to be false in sandbox mode, even when explicitly set to true")
	}
}

func TestEngine_K8s(t *testing.T) {
	cfg := &Config{
		Command: "echo hello",
		Engine:  "k8s",
		K8sNamespace: "default",
		K8sImage:     "alpine:latest",
	}

	cmd, err := New(cfg)
	if err != nil {
		// No cluster or kubeconfig: skip
		if strings.Contains(err.Error(), "k8s") || strings.Contains(err.Error(), "config") || strings.Contains(err.Error(), "in-cluster") {
			t.Skipf("k8s not available, skipping: %v", err)
		}
		t.Fatalf("failed to create command: %v", err)
	}

	buf := &strings.Builder{}
	cmd.SetStdout(buf)

	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "k8s") || strings.Contains(err.Error(), "job") || strings.Contains(err.Error(), "pod") {
			t.Skipf("k8s execution failed (cluster may be unavailable): %v", err)
		}
		t.Fatalf("failed to run command: %v", err)
	}

	if v := buf.String(); !strings.Contains(v, "hello") {
		t.Errorf("expected stdout to contain %q, got %q", "hello", v)
	}
}

func TestEngine_Podman(t *testing.T) {
	cfg := &Config{
		Command: "echo hello",
		Engine:  "podman",
		Image:   "alpine:latest",
	}

	cmd, err := New(cfg)
	if err != nil {
		if strings.Contains(err.Error(), "podman") || strings.Contains(err.Error(), "connect") || strings.Contains(err.Error(), "socket") {
			t.Skipf("podman not available, skipping: %v", err)
		}
		t.Fatalf("failed to create command: %v", err)
	}

	buf := &strings.Builder{}
	cmd.SetStdout(buf)

	if err := cmd.Run(); err != nil {
		if strings.Contains(err.Error(), "podman") || strings.Contains(err.Error(), "container") {
			t.Skipf("podman execution failed (podman may be unavailable): %v", err)
		}
		t.Fatalf("failed to run command: %v", err)
	}

	if v := buf.String(); !strings.Contains(v, "hello") {
		t.Errorf("expected stdout to contain %q, got %q", "hello", v)
	}
}

func TestEngine_WSL(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("wsl engine is only available on Windows")
	}
	cfg := &Config{
		Command: "echo hello",
		Engine:  "wsl",
	}
	cmd, err := New(cfg)
	if err != nil {
		if strings.Contains(err.Error(), "wsl") || strings.Contains(err.Error(), "Windows") {
			t.Skipf("wsl not available, skipping: %v", err)
		}
		t.Fatalf("failed to create command: %v", err)
	}
	buf := &strings.Builder{}
	cmd.SetStdout(buf)
	if err := cmd.Run(); err != nil {
		t.Skipf("wsl execution failed: %v", err)
	}
	if v := buf.String(); !strings.Contains(v, "hello") {
		t.Errorf("expected stdout to contain %q, got %q", "hello", v)
	}
}
