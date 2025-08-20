package main

import (
	"fmt"
	"net"
)

type Server struct {
	Addr    string
	Network string
	Router  Router
}

func main() {
	server := Server{
		Addr:    "localhost:8080",
		Network: "tcp",
		Router:  Router{},
	}

	fmt.Printf("Starting server")

	l, err := net.Listen(server.Network, server.Addr)
	if err != nil {
		panic(err)
	}

	defer l.Close()

	server.setupHandlers()

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading: %#v\n", err)
				return
			}

			fmt.Print("JSEM ZDE\n\n\n\n")
			fmt.Printf("Message received: %s\n", string(buf[:len]))
			r, err := parseRequest(string(buf[:len]))
			if err != nil {
				fmt.Print("\n\nSelhalo parsovani requestu\n\n")
				fmt.Printf("%v", err)
				response := Response{
					StatusBadRequest,
					"HTTP/1.1",
					make(map[string]string),
					"",
				}
				response.Headers["Content-Type"] = "text/plain"
				conn.Write([]byte(formatResponse(&response)))
				conn.Close()
				return
			}

			fmt.Printf("%v       %v", r.Method, r.Path)
			handler, err := server.Router.FindHandler(r)
			if err != nil {
				response := Response{
					StatusNotFound,
					"HTTP/1.1",
					make(map[string]string),
					"",
				}
				response.Headers["Content-Type"] = "text/plain"
				conn.Write([]byte(formatResponse(&response)))
				conn.Close()
				return
			}

			fmt.Print("Dosel jsem az sem")
			response := handler(r)

			conn.Write([]byte(formatResponse(response)))
			conn.Close()
		}(conn)
	}
}

func (s *Server) setupHandlers() {
	s.Router.AddHandler(GET, "/user/123", userHandler)
	s.Router.AddHandler(GET, "", userHandler)
	s.Router.AddHandler(GET, "/", userHandler)
	s.Router.AddHandler(GET, "/user/{id}/posts", userPostsHandler)
}

func userPostsHandler(r *Request) *Response {
	s := ""
	for key, value := range r.PathParams {
		s += key + " " + value + " "
	}
	response := Response{
		StatusOK,
		"HTTP/1.1",
		make(map[string]string),
		"Nazdar vitaj user id" + s,
	}
	response.Headers["Content-Type"] = "text/plain"
	return &response
}

func userHandler(*Request) *Response {
	response := Response{
		StatusOK,
		"HTTP/1.1",
		make(map[string]string),
		"Nazdar vitaj",
	}
	response.Headers["Content-Type"] = "text/plain"
	return &response
}
