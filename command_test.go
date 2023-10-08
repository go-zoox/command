package command

import (
	"context"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	cfg := &Config{
		Command: "echo hello world",
	}

	cmd, err := New(context.Background(), cfg)
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
