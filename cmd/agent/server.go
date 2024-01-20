package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/command/agent/server"
)

func registerServerCommand(app *cli.MultipleProgram) {
	app.Register("server", &cli.Command{
		Name:  "server",
		Usage: "Run command server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "Command server port",
				Aliases: []string{"p"},
				EnvVars: []string{"PORT"},
				Value:   8080,
			},
		},
		Action: func(ctx *cli.Context) error {
			s, err := server.New(func(opt *server.Option) {
				opt.Port = ctx.Int("port")
			})
			if err != nil {
				return err
			}

			return s.Run()
		},
	})
}
