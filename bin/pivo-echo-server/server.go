package main

import (
	"fmt"
	"log"

	"gopkg.in/pivo.v2"
)

type Server struct {
	accept   chan pivo.Connector
	conns    map[pivo.Connector]bool
	disconn  chan pivo.Connector
	read     chan *pivo.Message
	shutdown chan bool
}

func NewServer() *Server {
	return &Server{
		accept:   make(chan pivo.Connector),
		conns:    map[pivo.Connector]bool{},
		disconn:  make(chan pivo.Connector),
		read:     make(chan *pivo.Message),
		shutdown: make(chan bool),
	}
}

func (s *Server) DisconnectAll() {
	for conn, _ := range s.conns {
		conn.Close(nil)
	}
}

func (s *Server) Start() {
	// Forever loop
	for {
		select {

		// Accept
		case c := <-s.accept:
			log.Println("New connection from",
				c.RemoteAddr().String())
			s.conns[c] = true
			msg := fmt.Sprintf(WelcomeBanner,
				c.RemoteAddr().String(), pivo.Version)
			c.Send(pivo.TextMessage(nil, msg))

		// Disconnect
		case c := <-s.disconn:
			log.Println("Lost connection from",
				c.RemoteAddr().String())
			delete(s.conns, c)

		// Read
		case msg := <-s.read:
			log.Println("Received text from",
				msg.From.RemoteAddr().String(), ":",
				string(msg.Data))
			// Echo the message as-is
			msg.From.Send(msg)

		// Shutting down
		case <-s.shutdown:
			return
		}
	}
}

func (s *Server) Stop() { close(s.shutdown) }
