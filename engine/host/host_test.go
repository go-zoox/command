package host

import (
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	if Name != "host" {
		t.Errorf("Name = %q, want host", Name)
	}
}

func TestNew_SimpleCommand(t *testing.T) {
	cfg := &Config{
		Command: "echo host engine test",
		Shell:   "/bin/sh",
	}
	eng, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	buf := &strings.Builder{}
	_ = eng.SetStdout(buf)
	if err := eng.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if err := eng.Wait(); err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if !strings.Contains(buf.String(), "host engine test") {
		t.Errorf("stdout = %q, want to contain 'host engine test'", buf.String())
	}
}

func TestNew_EmptyShellUsesDefault(t *testing.T) {
	cfg := &Config{
		Command: "echo ok",
		Shell:   "",
	}
	eng, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if cfg.Shell != "/bin/sh" {
		t.Errorf("expected default Shell /bin/sh, got %q", cfg.Shell)
	}
	_ = eng
}
