package server

import (
	"log"
	"net"
)

type Server struct {
	addr     string
	listener net.Listener
	handle   handleFunc
}

type handleFunc func(conn net.Conn)

func New(addr string, handler handleFunc) (*Server, error) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{addr: addr, listener: server, handle: handler}, nil
}

func (s *Server) Run() {
	log.Printf("Server is running, Listen at %s\n", s.addr)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("accept err: ", err.Error())
		}
		go s.handle(conn)
	}
}
