package main

import (
	"strconv"
)

type Response struct {
	Status  HTTPStatus
	Version string
	Headers map[string]string
	Body    string
}

type HTTPStatus int

func formatResponse(r *Response) string {
	result := r.Version + " " + strconv.Itoa(int(r.Status)) + " " + StatusText(r.Status) + "\r\n"

	if len(r.Headers) > 0 {
		for key, value := range r.Headers {
			result += key + ": " + value + "\r\n"
		}
	}
	result += "\r\n"

	result += r.Body

	return result
}

const (
	// 2xx Success
	StatusOK        HTTPStatus = 200
	StatusCreated   HTTPStatus = 201
	StatusAccepted  HTTPStatus = 202
	StatusNoContent HTTPStatus = 204

	// 3xx Redirection
	StatusMovedPermanently  HTTPStatus = 301
	StatusFound             HTTPStatus = 302
	StatusNotModified       HTTPStatus = 304
	StatusTemporaryRedirect HTTPStatus = 307
	StatusPermanentRedirect HTTPStatus = 308

	// 4xx Client Error
	StatusBadRequest                   HTTPStatus = 400
	StatusUnauthorized                 HTTPStatus = 401
	StatusForbidden                    HTTPStatus = 403
	StatusNotFound                     HTTPStatus = 404
	StatusMethodNotAllowed             HTTPStatus = 405
	StatusNotAcceptable                HTTPStatus = 406
	StatusRequestTimeout               HTTPStatus = 408
	StatusConflict                     HTTPStatus = 409
	StatusGone                         HTTPStatus = 410
	StatusLengthRequired               HTTPStatus = 411
	StatusPreconditionFailed           HTTPStatus = 412
	StatusRequestEntityTooLarge        HTTPStatus = 413
	StatusRequestURITooLong            HTTPStatus = 414
	StatusUnsupportedMediaType         HTTPStatus = 415
	StatusRequestedRangeNotSatisfiable HTTPStatus = 416
	StatusExpectationFailed            HTTPStatus = 417
	StatusTeapot                       HTTPStatus = 418
	StatusTooManyRequests              HTTPStatus = 429

	// 5xx Server Error
	StatusInternalServerError     HTTPStatus = 500
	StatusNotImplemented          HTTPStatus = 501
	StatusBadGateway              HTTPStatus = 502
	StatusServiceUnavailable      HTTPStatus = 503
	StatusGatewayTimeout          HTTPStatus = 504
	StatusHTTPVersionNotSupported HTTPStatus = 505
)

func StatusText(status HTTPStatus) string {
	switch status {
	// 2xx Success
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusAccepted:
		return "Accepted"
	case StatusNoContent:
		return "No Content"

	// 3xx Redirection
	case StatusMovedPermanently:
		return "Moved Permanently"
	case StatusFound:
		return "Found"
	case StatusNotModified:
		return "Not Modified"
	case StatusTemporaryRedirect:
		return "Temporary Redirect"
	case StatusPermanentRedirect:
		return "Permanent Redirect"

	// 4xx Client Error
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusNotAcceptable:
		return "Not Acceptable"
	case StatusRequestTimeout:
		return "Request Timeout"
	case StatusConflict:
		return "Conflict"
	case StatusGone:
		return "Gone"
	case StatusLengthRequired:
		return "Length Required"
	case StatusPreconditionFailed:
		return "Precondition Failed"
	case StatusRequestEntityTooLarge:
		return "Request Entity Too Large"
	case StatusRequestURITooLong:
		return "Request-URI Too Long"
	case StatusUnsupportedMediaType:
		return "Unsupported Media Type"
	case StatusRequestedRangeNotSatisfiable:
		return "Requested Range Not Satisfiable"
	case StatusExpectationFailed:
		return "Expectation Failed"
	case StatusTeapot:
		return "I'm a teapot"
	case StatusTooManyRequests:
		return "Too Many Requests"

	// 5xx Server Error
	case StatusInternalServerError:
		return "Internal Server Error"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	case StatusGatewayTimeout:
		return "Gateway Timeout"
	case StatusHTTPVersionNotSupported:
		return "HTTP Version Not Supported"

	default:
		return "Unknown Status"
	}
}
