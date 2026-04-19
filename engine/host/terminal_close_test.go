package host

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/creack/pty"
)

func TestTerminal_Close_Idempotent(t *testing.T) {
	// Match a typical PTY session: shell waits on a child (same pattern as production).
	cmd := exec.Command("/bin/sh", "-c", "sleep 30 & wait $!")
	f, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("pty.Start: %v", err)
	}
	term := &Terminal{File: f, Cmd: cmd}

	if err := term.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := term.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
	_ = cmd.Wait()
}

// TestTerminal_Close_KillsShellChildProcessGroup verifies Close tears down a long-running
// child of /bin/sh -c (same process group as the session leader from pty.Start).
func TestTerminal_Close_KillsShellChildProcessGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("spawns sleep subprocess")
	}

	dir := t.TempDir()
	pidfile := filepath.Join(dir, "child.pid")
	// Background sleep, record PID, then wait so the shell stays alive (otherwise an
	// exiting login shell may SIGHUP the background job before we call Close).
	script := fmt.Sprintf("sleep 86400 & echo $! > '%s'; wait $!", pidfile)
	cmd := exec.Command("/bin/sh", "-c", script)
	f, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("pty.Start: %v", err)
	}
	term := &Terminal{File: f, Cmd: cmd}

	var childPID int
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		b, err := os.ReadFile(pidfile)
		if err == nil && len(b) > 0 {
			childPID, err = strconv.Atoi(strings.TrimSpace(string(b)))
			if err == nil && childPID > 0 {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	if childPID <= 0 {
		t.Fatal("timed out waiting for child PID file")
	}

	if err := syscall.Kill(childPID, 0); err != nil {
		t.Fatalf("child %d not running before Close: %v", childPID, err)
	}

	if err := term.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	_ = cmd.Wait()

	time.Sleep(50 * time.Millisecond)
	if err := syscall.Kill(childPID, 0); err == nil {
		t.Fatalf("child sleep %d still alive after Close", childPID)
	} else if errno, ok := err.(syscall.Errno); !ok || errno != syscall.ESRCH {
		t.Fatalf("expected ESRCH for dead child, got %v", err)
	}
}
