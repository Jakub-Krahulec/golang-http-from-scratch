package main

import (
	"fmt"
	"slices"
)

type HandlerFunc func(r *Request) *Response
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(r *Request) *Response {
		fmt.Printf("METHOD: %v, PATH: %v, HTTP VERSION: %v, BODY: %v\n", r.Method, r.Path, r.Version, r.Body)
		return next(r)
	}
}

func RecoveryMiddleware(next HandlerFunc) HandlerFunc {
	return func(r *Request) (response *Response) {
		defer func() {
			if r := recover(); r != nil {
				response = NewResponse(StatusInternalServerError, PlainTextHeaders(), "")
			}
		}()

		response = next(r)
		return
	}
}

// Cross Origin Resource Sharing
func corsMiddleware(next HandlerFunc) HandlerFunc {
	return func(r *Request) *Response {
		// maybe move into helper method later if needed
		allowedOrigins := []string{}
		origin := r.Headers["origin"]
		if origin == "" {
			return next(r)
		}
		if slices.Contains(allowedOrigins, origin) {
			response := next(r)
			// move into helper function if needed later
			if response.Headers == nil {
				response.Headers = make(map[string]string)
			}
			response.Headers["Access-Control-Allow-Origin"] = origin
			return response
		}

		fmt.Print("CORS policy violation\n")
		return NewResponse(StatusForbidden, PlainTextHeaders(), "")
	}
}
