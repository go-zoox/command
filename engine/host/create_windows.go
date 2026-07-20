//go:build windows

package host

import (
	"os/exec"
	"os/user"

	"github.com/go-zoox/logger"
)

func applyUser(cmd *exec.Cmd, username string) error {
	if username == "" {
		return nil
	}

	userX, err := user.Lookup(username)
	if err != nil {
		return err
	}

	logger.Infof("[command] uid=%s gid=%s (windows: user switching not fully supported)", userX.Uid, userX.Gid)

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
