package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/command"
	"github.com/go-zoox/command/terminal"

	"golang.org/x/term"
)

func Exec(app *cli.MultipleProgram) {
	app.Register("exec", &cli.Command{
		Name:  "exec",
		Usage: "command execute",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "engine",
				Usage:   "command engine",
				Aliases: []string{"e"},
				EnvVars: []string{"ENGINE"},
				Value:   "host",
			},
			&cli.StringFlag{
				Name:    "command",
				Usage:   "the command",
				Aliases: []string{"c"},
				EnvVars: []string{"COMMAND"},
			},
			&cli.StringFlag{
				Name:    "workdir",
				Usage:   "the command workdir",
				Aliases: []string{"w"},
				EnvVars: []string{"WORKDIR"},
			},
			&cli.StringFlag{
				Name:    "user",
				Usage:   "the command user",
				Aliases: []string{"u"},
				// EnvVars: []string{"WORKDIR"},
			},
			&cli.StringFlag{
				Name:    "shell",
				Usage:   "the command shell",
				Aliases: []string{"s"},
				// EnvVars: []string{"SHELL"},
			},
			&cli.StringFlag{
				Name:    "image",
				Usage:   "docker image",
				Aliases: []string{"i"},
				EnvVars: []string{"IMAGE"},
			},
			&cli.BoolFlag{
				Name:    "tty",
				Usage:   "Allocate a pseudo-TTY. The default is false, which disables TTY allocation.",
				Aliases: []string{"t"},
				EnvVars: []string{"TTY"},
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			cmd, err := command.New(context.Background(), &command.Config{
				Engine:  ctx.String("engine"),
				Command: ctx.String("command"),
				WorkDir: ctx.String("workdir"),
				User:    ctx.String("user"),
				Shell:   ctx.String("shell"),
				Image:   ctx.String("image"),
			})
			if err != nil {
				return err
			}

			if ctx.Bool("tty") {
				term, err := cmd.Terminal()
				if err != nil {
					return err
				}
				defer term.Close()

				go func() {
					io.Copy(os.Stdout, term)
					// _, err := io.Copy(os.Stdout, term)
					// if err != nil && err != io.EOF {
					// 	os.Stderr.Write([]byte(err.Error()))
					// 	os.Exit(term.ExitCode())
					// 	return
					// }
				}()

				if err := connectKeyboard(term); err != nil {
					return err
				}

				return nil
			}

			return cmd.Run()
		},
	})
}

func connectKeyboard(t terminal.Terminal) error {
	// resize
	if err := resizeTerminal(t); err != nil {
		return err
	}

	// 监听操作系统信号
	sigWinch := make(chan os.Signal, 1)
	signal.Notify(sigWinch, syscall.SIGWINCH)
	// 启动循环来检测终端窗口大小是否发生变化
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

	if err := keyboard.Open(); err != nil {
		return err
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		// fmt.Printf("You pressed: rune:%q, key %X\r\n", char, key)
		if key == keyboard.KeyCtrlD {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// key == 0 => char
		if key == 0 {
			_, err = t.Write([]byte{byte(char)})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			switch key {
			case keyboard.KeyF1:
				_, err = t.Write([]byte{0x1b, 0x4f, 0x50})
			case keyboard.KeyF2:
				_, err = t.Write([]byte{0x1b, 0x4f, 0x51})
			case keyboard.KeyF3:
				_, err = t.Write([]byte{0x1b, 0x4f, 0x52})
			case keyboard.KeyF4:
				_, err = t.Write([]byte{0x1b, 0x4f, 0x53})
			case keyboard.KeyF5:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x31, 0x35, 0x7e})
			case keyboard.KeyF6:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x31, 0x37, 0x7e})
			case keyboard.KeyF7:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x31, 0x38, 0x7e})
			case keyboard.KeyF8:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x31, 0x39, 0x7e})
			case keyboard.KeyF9:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x32, 0x30, 0x7e})
			case keyboard.KeyF10:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x32, 0x31, 0x7e})
			case keyboard.KeyF11:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x32, 0x33, 0x7e})
			case keyboard.KeyF12:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x32, 0x34, 0x7e})
			case keyboard.KeyInsert:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x32, 0x7e})
			case keyboard.KeyDelete:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x33, 0x7e})
			case keyboard.KeyHome:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x48})
			case keyboard.KeyEnd:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x46})
			case keyboard.KeyPgup:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x35, 0x7e})
			case keyboard.KeyPgdn:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x36, 0x7e})
			case keyboard.KeyArrowUp:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x41})
			case keyboard.KeyArrowDown:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x42})
			case keyboard.KeyArrowRight:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x43})
			case keyboard.KeyArrowLeft:
				_, err = t.Write([]byte{0x1b, 0x5b, 0x44})
			default:
				_, err = t.Write([]byte{byte(key)})
			}

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}

	return nil
}

func resizeTerminal(t terminal.Terminal) error {
	fd := int(os.Stdin.Fd())
	columns, rows, err := term.GetSize(fd)
	if err != nil {
		return err
	}

	return t.Resize(rows, columns)
}
