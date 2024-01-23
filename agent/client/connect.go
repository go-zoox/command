package client

import (
	"github.com/go-zoox/command/agent/event"
	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/core-utils/fmt"
	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket"
	"github.com/go-zoox/websocket/conn"
)

func (c *client) Connect() error {
	logger.Debugf("create websocket client to %s", c.opt.Server)

	if c.opt.Server == "" {
		return fmt.Errorf("server address is required")
	}

	ws, err := websocket.NewClient(func(opt *websocket.ClientOption) {
		opt.Addr = c.opt.Server
	})
	if err != nil {
		return err
	}

	connected := make(chan struct{})

	logger.Debugf("listen on close event ...")
	ws.OnClose(func(conn conn.Conn, code int, message string) error {
		c.stderr.Write([]byte(message))
		c.exitcodeCh <- code
		return nil
	})

	logger.Debugf("listen on connect event ...")
	ws.OnConnect(func(cc websocket.Conn) error {
		c.core = cc

		connected <- struct{}{}

		cc.OnMessage(func(typ int, message []byte) error {
			if typ != conn.TextMessage {
				return nil
			}

			handleMessage := func() error {
				evt := &event.Event{}
				if err := evt.Decode(message); err != nil {
					return err
				}

				logger.Debugf("receive event from server: %s", message)

				switch evt.Type {
				case event.Done:
					responseEvent := &event.DoneEvent{}
					if err := responseEvent.Decode(message); err != nil {
						return err
					}
					switch string(responseEvent.Payload) {
					case event.New:
						c.newEventDone <- struct{}{}
					case event.Start:
						c.startEventDone <- struct{}{}
					case event.Wait:
						c.waitEventDone <- struct{}{}
					case event.Cancel:
						c.cancelEventDone <- struct{}{}
					}
				case event.Stdout:
					stdoutEvent := &event.StdoutEvent{}
					if err := stdoutEvent.Decode(message); err != nil {
						return err
					}

					c.stdout.Write(stdoutEvent.Payload)
				case event.Stderr:
					stderrEvent := &event.StderrEvent{}
					if err := stderrEvent.Decode(message); err != nil {
						return err
					}

					c.stderr.Write(stderrEvent.Payload)
				case event.Exitcode:
					exitcodeEvent := &event.ExitcodeEvent{}
					if err := exitcodeEvent.Decode(message); err != nil {
						return err
					}

					exitcode := cast.ToInt(string(exitcodeEvent.Payload))
					c.exitcodeCh <- exitcode
				}
				return nil
			}

			go func() {
				if err := handleMessage(); err != nil {
					logger.Errorf("handle message error: %s", err)
				}
			}()

			return nil
		})

		return nil
	})

	logger.Debugf("connect to server ...")
	if err := ws.Connect(); err != nil {
		return err
	}

	<-connected
	logger.Debugf("connected")
	return nil
}
