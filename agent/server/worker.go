package server

import (
	"encoding/json"
	"fmt"
	gio "io"

	"github.com/go-zoox/command"
	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/command/errors"
	"github.com/go-zoox/core-utils/io"
	"github.com/go-zoox/eventemitter"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
)

func Worker(c conn.Conn) {
	// @2 utils
	sendEvent := func(evt *event.Event) error {
		s, err := json.Marshal(evt)
		if err != nil {
			return err
		}

		return c.WriteTextMessage(s)
	}

	createWriter := func(eventName string) gio.Writer {
		return io.WriterWrapFunc(func(b []byte) (n int, err error) {
			err = sendEvent(&event.Event{
				Type:    eventName,
				Payload: b,
			})
			if err != nil {
				return 0, err
			}

			return len(b), nil
		})
	}

	// @1 vars
	eventBus := eventemitter.New(func(opt *eventemitter.Option) {
		opt.Context = c.Context()
	})
	//
	stdout := createWriter(event.Stdout)
	stderr := createWriter(event.Stderr)
	exitcode := createWriter(event.Exitcode)
	//
	done := createWriter(event.Done)
	//
	var cmd command.Command
	var cfg *command.Config

	// @3 connection listeners
	c.OnMessage(func(typ int, message []byte) error {
		if typ != conn.TextMessage {
			return nil
		}

		evt := &event.Event{}
		if err := evt.Decode(message); err != nil {
			return err
		}

		switch evt.Type {
		case event.New:
			newEvent := &event.NewEvent{}
			if err := newEvent.Decode(message); err != nil {
				return err
			}
			eventBus.Emit(event.New, newEvent.Payload)
		case event.Start:
			eventBus.Emit(event.Start, nil)
		case event.Wait:
			eventBus.Emit(event.Wait, nil)
		case event.Cancel:
			eventBus.Emit(event.Cancel, nil)
		//
		case event.Stdin:
			stdinEvent := &event.StdinEvent{}
			if err := stdinEvent.Decode(message); err != nil {
				return err
			}

			eventBus.Emit(event.Stdin, stdinEvent.Payload)
		default:
			return fmt.Errorf("unknown event type: %s", evt.Type)
		}

		return nil
	})

	// @4 event listeners
	eventBus.On("error", eventemitter.HandleFunc(func(payload any) {
		err, ok := payload.(error)
		if !ok {
			return
		}

		if errx, ok := err.(*errors.ExitError); ok {
			logger.Debugf("failed to run command: %s (exit code: %d)", cfg.Command, errx.ExitCode())
			exitcode.Write([]byte(fmt.Sprintf("%d", errx.ExitCode())))
			return
		}

		logger.Debugf("failed to run command(2): %s (exit code: %d)", cfg.Command, 1)
		exitcode.Write([]byte("1"))
	}))

	eventBus.On(event.New, eventemitter.HandleFunc(func(payload any) {
		defer done.Write([]byte(event.New))

		logger.Debugf("[stage:%s] create command ...", event.New)
		if cmd != nil {
			eventBus.Emit("error", fmt.Errorf("[stage:%s] command is already created", event.New))
			return
		}

		cfg = payload.(*command.Config)
		cfg.Context = c.Context()

		cm, err := command.New(cfg)
		if err != nil {
			eventBus.Emit("error", err)
			return
		}
		cmd = cm

		// cmd.SetStdin(stdin)
		cmd.SetStdout(stdout)
		cmd.SetStderr(stderr)
	}))

	eventBus.On(event.Start, eventemitter.HandleFunc(func(payload any) {
		defer done.Write([]byte(event.Start))

		logger.Debugf("[stage:%s] start command ...", event.Start)
		if cmd == nil {
			eventBus.Emit("error", fmt.Errorf("[stage:%s] command is not created", event.Start))
			return
		}

		if err := cmd.Start(); err != nil {
			eventBus.Emit("error", err)
			return
		}

	}))

	eventBus.On(event.Wait, eventemitter.HandleFunc(func(payload any) {
		defer done.Write([]byte(event.Wait))

		logger.Debugf("[stage:%s] wait for command ...", event.Wait)
		if cmd == nil {
			eventBus.Emit("error", fmt.Errorf("[stage:%s]  command is not created", event.Wait))
			return
		}

		if err := cmd.Wait(); err != nil {
			eventBus.Emit("error", err)
			return
		}

		logger.Debugf("[stage:%s] command is done", event.Wait)
		exitcode.Write([]byte("0"))
	}))

	eventBus.On(event.Cancel, eventemitter.HandleFunc(func(payload any) {
		done.Write([]byte(event.Cancel))

		logger.Debugf("[stage:%s] cancel command ...", event.Cancel)
		if cmd == nil {
			eventBus.Emit("error", fmt.Errorf("[stage:%s]  command is not created", event.Cancel))
			return
		}

		if err := cmd.Cancel(); err != nil {
			eventBus.Emit("error", err)
			return
		}
	}))
}
