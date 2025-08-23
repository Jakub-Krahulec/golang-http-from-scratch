package main

import (
	"fmt"
)

type HandlerFunc func(r *Request) *Response
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(r *Request) *Response {
		fmt.Printf("Zde budu neco logovat")
		return next(r)
	}
}
