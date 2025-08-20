package main

type Request struct {
	Method       HTTPMethod
	Path         string
	PathSegments []string
	Version      string
	Headers      map[string]string
	Body         string
	QueryParams  map[string]string
	PathParams   map[string]string
}

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	PATCH  HTTPMethod = "PATCH"
	DELETE HTTPMethod = "DELETE"
)
