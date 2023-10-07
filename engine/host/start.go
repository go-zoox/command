package host

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"syscall"

	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/logger"
)

func (h *host) Start() error {
	if h.cmd != nil {
		return errors.New("command: already started")
	}

	logger.Debugf("start command: %s %s %s", h.cfg.Shell, "-c", h.cfg.Command)
	h.cmd = exec.Command(h.cfg.Shell, "-c", h.cfg.Command)

	if err := applyEnv(h.cmd, h.cfg.Environment); err != nil {
		return err
	}

	if err := applyWorkDir(h.cmd, h.cfg.WorkDir); err != nil {
		return err
	}

	if err := applyUser(h.cmd, h.cfg.User); err != nil {
		return err
	}

	if err := applyHistory(h.cmd, h.cfg.IsHistoryDisabled); err != nil {
		return err
	}

	if err := applyStdin(h.cmd, h.stdin); err != nil {
		return nil
	}

	if err := applyStdout(h.cmd, h.stdout); err != nil {
		return nil
	}

	if err := applyStderr(h.cmd, h.stderr); err != nil {
		return nil
	}

	return h.cmd.Start()
}

func applyEnv(cmd *exec.Cmd, environment map[string]string) error {
	cmd.Env = append(os.Environ(), "TERM=xterm")

	for k, v := range environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return nil
}

func applyWorkDir(cmd *exec.Cmd, workDir string) error {
	cmd.Dir = workDir
	return nil
}

func applyUser(cmd *exec.Cmd, username string) error {
	if username == "" {
		return nil
	}

	userX, err := user.Lookup(username)
	if err != nil {
		return err
	}

	logger.Infof("[command] uid=%s gid=%s", userX.Uid, userX.Gid)

	uid := cast.ToInt(userX.Uid)
	gid := cast.ToInt(userX.Gid)

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(uid),
		Gid: uint32(gid),
	}

	cmd.Env = append(
		cmd.Env,
		"USER="+username,
		"HOME="+userX.HomeDir,
		"LOGNAME="+username,
		"UID="+userX.Uid,
		"GID="+userX.Gid,
	)

	return nil
}

func applyHistory(cmd *exec.Cmd, disable bool) error {
	if disable {
		cmd.Env = append(cmd.Env, "HISTFILE=/dev/null")
	}

	return nil
}

func applyStdin(cmd *exec.Cmd, stdin io.Reader) error {
	cmd.Stdin = stdin
	return nil
}

func applyStdout(cmd *exec.Cmd, stdout io.Writer) error {
	cmd.Stdout = stdout
	if cmd.Stderr == nil {
		return applyStderr(cmd, stdout)
	}

	return nil
}

func applyStderr(cmd *exec.Cmd, stderr io.Writer) error {
	cmd.Stderr = stderr
	return nil
}
