package main

func main() {
	server := setupServer()
	server.startListening()
	defer server.Listener.Close()
	server.setupHandlers()
	server.handleConnections()
}

func setupServer() Server {
	server := Server{
		Addr:              "localhost:8080",
		Network:           "tcp",
		Router:            Router{},
		GlobalMiddlewares: []MiddlewareFunc{corsMiddleware, LoggingMiddleware, RecoveryMiddleware},
	}
	return server
}
