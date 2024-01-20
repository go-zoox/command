package server

import (
	"fmt"

	"github.com/go-zoox/websocket"
	"github.com/go-zoox/websocket/conn"
)

func (s *server) Run() error {
	ws, err := websocket.NewServer()
	if err != nil {
		return err
	}

	ws.OnClose(func(conn conn.Conn, code int, message string) error {
		return nil
	})

	ws.OnError(func(conn conn.Conn, err error) error {
		return nil
	})

	ws.OnConnect(func(conn conn.Conn) error {
		go Worker(conn)

		return nil
	})

	return ws.Run(fmt.Sprintf(":%d", s.opt.Port))
}
