package main

import (
	"errors"
	"log"
	"net/http"

	"gopkg.in/pivo.v2/ws"
)

const websocketEchoUri = `/`

var ErrProtocolViolation = errors.New("protocol violation")

type websocket struct {
	conn *ws.Conn
}

func (ws *websocket) OnClose(why error) error {
	server.disconn <- ws.conn
	return nil
}

func (ws *websocket) OnBinaryRead(data []byte) error {
	// Binary transfer is not allowed here
	return ErrProtocolViolation
}

func (ws *websocket) OnTextRead(text string) error {
	server.read <- ws.conn.TextMessage(text)
	return nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := ws.DefaultConn()
	if err := conn.Upgrade(w, r, nil); err != nil {
		log.Printf("websocket: %s: failed to upgrade: %s",
			r.RemoteAddr, err)
		return
	}
	go conn.Sender()
	go conn.Receiver(&websocket{conn})
	server.accept <- conn
}

func init() {
	http.Handle(websocketEchoUri, server)
}
