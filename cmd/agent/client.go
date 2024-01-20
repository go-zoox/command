package main

import (
	"os"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/command"
	"github.com/go-zoox/command/agent/client"
	"github.com/go-zoox/command/errors"
)

func registerClientCommand(app *cli.MultipleProgram) {
	app.Register("client", &cli.Command{
		Name:  "client",
		Usage: "Run command agent client",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "command",
				Usage:    "Command to run",
				Aliases:  []string{"c"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "server",
				Usage:   "Command server address",
				Aliases: []string{"s"},
				EnvVars: []string{"SERVER"},
				Value:   "ws://localhost:8080",
			},
		},
		Action: func(ctx *cli.Context) error {
			s := client.New(func(opt *client.Option) {
				opt.Server = ctx.String("server")
			})

			if err := s.Connect(); err != nil {
				return err
			}

			err := s.New(&command.Config{
				Command: ctx.String("command"),
				// Timeout: 1 * time.Microsecond,
			})
			if err != nil {
				if errx, ok := err.(*errors.ExitError); ok {
					os.Exit(errx.ExitCode())
					return nil
				}

				return err
			}

			err = s.Run()
			if err != nil {
				if errx, ok := err.(*errors.ExitError); ok {
					os.Exit(errx.ExitCode())
					return nil
				}
			}

			return err
		},
	})
}
