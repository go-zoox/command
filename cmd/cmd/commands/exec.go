package commands

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-zoox/core-utils/strings"

	"github.com/eiannone/keyboard"
	"github.com/go-zoox/cli"
	"github.com/go-zoox/command"
	"github.com/go-zoox/command/terminal"

	"golang.org/x/term"
)

// Exec is the exec command
func Exec(app *cli.MultipleProgram) {
	app.Register("exec", &cli.Command{
		Name:  "exec",
		Usage: "command execute",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "agent",
				Usage:   "command agent server",
				EnvVars: []string{"AGENT"},
			},
			&cli.StringFlag{
				Name:    "engine",
				Usage:   "command engine, avaliable: host, docker, caas",
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
				// Value:   fs.CurrentDir(),
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
			&cli.Int64Flag{
				Name:    "memory",
				Usage:   `Memory limit, unit: MB`,
				Aliases: []string{"m"},
				EnvVars: []string{"MEMORY"},
			},
			&cli.Float64Flag{
				Name:    "cpu",
				Usage:   `CPU limit, unit: core`,
				EnvVars: []string{"CPU"},
			},
			&cli.StringFlag{
				Name:    "platform",
				Usage:   `Command platform, available: linux/amd64, linux/arm64`,
				Aliases: []string{"p"},
				EnvVars: []string{"PLATFORM"},
			},
			&cli.StringFlag{
				Name:    "network",
				Usage:   `Network name`,
				Aliases: []string{"n"},
				EnvVars: []string{"NETWORK"},
			},
			&cli.BoolFlag{
				Name:    "disable-network",
				Usage:   "Disable network visit",
				EnvVars: []string{"DISABLE_NETWORK"},
			},
			&cli.StringFlag{
				Name:    "server",
				Usage:   "Command server address",
				EnvVars: []string{"SERVER"},
			},
			&cli.StringFlag{
				Name:    "client-id",
				Usage:   `Client ID for Server Auth`,
				EnvVars: []string{"CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:    "client-secret",
				Usage:   `Client Secret for Server Auth`,
				EnvVars: []string{"CLIENT_SECRET"},
			},
			// ssh
			&cli.StringFlag{
				Name:    "ssh-host",
				Usage:   "SSH host",
				EnvVars: []string{"SSH_HOST"},
			},
			&cli.IntFlag{
				Name:    "ssh-port",
				Usage:   "SSH port",
				EnvVars: []string{"SSH_PORT"},
				Value:   22,
			},
			&cli.StringFlag{
				Name:    "ssh-user",
				Usage:   "SSH user",
				EnvVars: []string{"SSH_USER"},
			},
			&cli.StringFlag{
				Name:    "ssh-pass",
				Usage:   "SSH password",
				EnvVars: []string{"SSH_PASS"},
			},
			&cli.StringFlag{
				Name:    "ssh-private-key",
				Usage:   "SSH private key",
				EnvVars: []string{"SSH_PRIVATE_KEY"},
			},
			&cli.StringFlag{
				Name:    "ssh-private-key-secret",
				Usage:   "SSH private key secret",
				EnvVars: []string{"SSH_PRIVATE_KEY_SECRET"},
			},
			&cli.BoolFlag{
				Name:    "ssh-is-ignore-strict-host-key-checking",
				Usage:   "SSH is ignore strict host key checking",
				EnvVars: []string{"SSH_IS_IGNORE_STRICT_HOST_KEY_CHECKING"},
			},
			&cli.StringFlag{
				Name:    "ssh-know-hosts-file-path",
				Usage:   "SSH know hosts file path",
				EnvVars: []string{"SSH_KNOW_HOSTS_FILE_PATH"},
			},
			&cli.StringSliceFlag{
				Name:  "env",
				Usage: `Set environment variables for the command`,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			environment := map[string]string{}
			for _, e := range ctx.StringSlice("env") {
				kv := strings.SplitN(e, "=", 2)
				environment[kv[0]] = kv[1]
			}

			cmd, err := command.New(&command.Config{
				Agent: ctx.String("agent"),
				//
				Engine:         ctx.String("engine"),
				Command:        ctx.String("command"),
				Environment:    environment,
				WorkDir:        ctx.String("workdir"),
				User:           ctx.String("user"),
				Shell:          ctx.String("shell"),
				Image:          ctx.String("image"),
				Memory:         ctx.Int64("memory"),
				CPU:            ctx.Float64("cpu"),
				Platform:       ctx.String("platform"),
				Network:        ctx.String("network"),
				DisableNetwork: ctx.Bool("disable-network"),
				//
				Server:       ctx.String("server"),
				ClientID:     ctx.String("client-id"),
				ClientSecret: ctx.String("client-secret"),
				//
				SSHHost:                          ctx.String("ssh-host"),
				SSHPort:                          ctx.Int("ssh-port"),
				SSHUser:                          ctx.String("ssh-user"),
				SSHPass:                          ctx.String("ssh-pass"),
				SSHPrivateKey:                    ctx.String("ssh-private-key"),
				SSHPrivateKeySecret:              ctx.String("ssh-private-key-secret"),
				SSHIsIgnoreStrictHostKeyChecking: ctx.Bool("ssh-is-ignore-strict-host-key-checking"),
				SSHKnowHostsFilePath:             ctx.String("ssh-know-hosts-file-path"),
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
