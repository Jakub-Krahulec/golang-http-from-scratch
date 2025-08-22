package main

import (
	"errors"
	"strings"
)

type Router struct {
	Routes []Route
}

type Route struct {
	Path    string
	Method  HTTPMethod
	Handler func(*Request) *Response
}

func (r *Router) AddHandler(method HTTPMethod, path string, handler func(*Request) *Response) {
	r.Routes = append(r.Routes, Route{path, method, handler})
}

func (r *Router) FindHandler(req *Request) (func(*Request) *Response, error) {
	// prvni zkusim najit presnou cestu
	for _, route := range r.Routes {
		if route.Method == req.Method && route.Path == req.Path {
			return route.Handler, nil
		}
	}

	for _, route := range r.Routes {
		if route.Method != req.Method {
			continue
		}
		splitedRoute := splitHTTPPathIntoSegments(route.Path)

		if len(splitedRoute) != len(req.PathSegments) {
			continue
		}

		foundMatch := true
		for i := range splitedRoute {
			if strings.HasPrefix(splitedRoute[i], "{") && strings.HasSuffix(splitedRoute[i], "}") {
				if req.PathParams == nil {
					req.PathParams = make(map[string]string)
				}
				paramName := splitedRoute[i][1 : len(splitedRoute[i])-1]
				req.PathParams[paramName] = req.PathSegments[i]
				continue
			}
			if splitedRoute[i] != req.PathSegments[i] {
				foundMatch = false
				break
			}
		}
		if foundMatch {
			return route.Handler, nil
		}
	}

	return nil, errors.New("invalid route")
}
