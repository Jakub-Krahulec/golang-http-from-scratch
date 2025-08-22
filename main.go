package main

import (
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Addr     string
	Network  string
	Router   Router
	Listener net.Listener
}

func main() {
	server := setupServer()
	server.startListening()
	defer server.Listener.Close()
	server.setupHandlers()
	server.handleConnections()
}

func setupServer() Server {
	server := Server{
		Addr:    "localhost:8080",
		Network: "tcp",
		Router:  Router{},
	}
	return server
}

func (s *Server) startListening() {
	fmt.Printf("Starting server")

	l, err := net.Listen(s.Network, s.Addr)
	if err != nil {
		panic(err)
	}
	s.Listener = l
}

func (s *Server) handleConnections() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()
			buf, bufLen, err := readBytes(conn)
			if err != nil {
				return
			}
			r, err := parseRequest(string(buf[:bufLen]))
			if err != nil {
				response := NewResponse(StatusBadRequest, PlainTextHeaders(), "Could not parse request")
				conn.Write([]byte(formatResponse(response)))
				return
			}

			handler, err := s.Router.FindHandler(r)
			if err != nil {
				response := NewResponse(StatusNotFound, PlainTextHeaders(), "")
				conn.Write([]byte(formatResponse(response)))
				return
			}

			response := handler(r)

			conn.Write([]byte(formatResponse(response)))
		}(conn)
	}
}

func readBytes(conn net.Conn) ([]byte, int, error) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading: %#v\n", err)
		return make([]byte, 1), 0, errors.New("cannot read connection")
	}
	return buf, len, nil
}

func (s *Server) setupHandlers() {
	s.Router.AddHandler(GET, "/user/123", userHandler)
	s.Router.AddHandler(GET, "", userHandler)
	s.Router.AddHandler(GET, "/", userHandler)
	s.Router.AddHandler(GET, "/user/{id}/posts", userPostsHandler)
}
