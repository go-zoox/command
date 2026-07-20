//go:build windows

package host

import (
	"os/exec"
	"os/user"
	"testing"
)

func TestApplyUser_EmptyUsername(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "echo hello")
	err := applyUser(cmd, "")
	if err != nil {
		t.Fatalf("applyUser with empty username should not error: %v", err)
	}
}

func TestApplyUser_NonexistentUser(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "echo hello")
	err := applyUser(cmd, "nonexistent_user_abcdef123456789")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestApplyUser_CurrentUser(t *testing.T) {
	currentUser, err := user.Current()
	if err != nil {
		t.Skipf("cannot get current user: %v", err)
	}

	cmd := exec.Command("cmd", "/c", "echo hello")
	err = applyUser(cmd, currentUser.Username)
	if err != nil {
		t.Fatalf("applyUser for current user should not error: %v", err)
	}

	foundUser := false
	foundHome := false
	for _, env := range cmd.Env {
		if env == "USER="+currentUser.Username {
			foundUser = true
		}
		if env == "HOME="+currentUser.HomeDir {
			foundHome = true
		}
	}
	if !foundUser {
		t.Error("USER env var not set in command environment")
	}
	if !foundHome {
		t.Error("HOME env var not set in command environment")
	}
}

func TestKillProcess_Success(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "ping -n 10 127.0.0.1 > nul")
	if err := cmd.Start(); err != nil {
		t.Fatalf("could not start process: %v", err)
	}
	if cmd.Process == nil {
		t.Fatal("process should not be nil after start")
	}

	err := killProcess(cmd.Process)
	if err != nil {
		t.Errorf("killProcess should not error on valid process: %v", err)
	}

	_ = cmd.Wait()
}

func TestKillProcess_AlreadyDead(t *testing.T) {
	cmd := exec.Command("cmd", "/c", "exit 0")
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not run process: %v", err)
	}

	err := killProcess(cmd.Process)
	if err != nil {
		t.Logf("killProcess on dead process returned: %v (may be expected)", err)
	}
}
