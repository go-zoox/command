package command

import (
	"strings"
	"testing"
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
