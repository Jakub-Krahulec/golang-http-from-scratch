package main

import (
	"fmt"
)

type HandlerFunc func(r *Request) *Response
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(r *Request) *Response {
		fmt.Printf("METHOD: %v, PATH: %v, HTTP VERSION: %v, BODY: %v", r.Method, r.Path, r.Version, r.Body)
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
