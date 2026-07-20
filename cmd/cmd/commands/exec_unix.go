//go:build !windows

package commands

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-zoox/command/terminal"
)

func startResizeWatcher(t terminal.Terminal) {
	sigWinch := make(chan os.Signal, 1)
	signal.Notify(sigWinch, syscall.SIGWINCH)
	go func() {
		for {
			select {
			case <-sigWinch:
				resizeTerminal(t)
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
}
