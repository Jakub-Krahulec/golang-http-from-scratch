package main

import (
	"encoding/hex"
	"errors"
	"strings"
)

// TODO: Finish secturity after reading about it and move checks to right places
func checkSecurityIssues(s string) error {
	// CRLF Injection -
	// CR and LF are special characters in ASCII 13 and 10 - /r /n1 - used in windows and internet protocols including http
	// it is used to add http headers into http response
	// also used to falsifie logs by inserting lines and hide antoher attacks or confuse admins
	if strings.Contains(s, "\r\r") || strings.Contains(s, "\n\n") {
		return errors.New("Potential CRLF injection")
	}

	// Path traversal protection
	// used to access files outside
	if strings.Contains(s, "..") {
		return errors.New("Potential Path traversal protection")
	}

	// Null byte injection
	// used in C to detect end of the string
	// bypass file extension check,
	if strings.Contains(s, "\x00") {
		return errors.New("Potential null byte injection")
	}

	return nil
}

func parseRequest(s string) (*Request, error) {
	lines := strings.Split(s, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("could note pars request (1)")
	}

	err := checkSecurityIssues(s)
	if err != nil {
		return nil, err
	}

	var r Request
	err = parseFirstLine(lines[0], &r)
	if err != nil {
		return nil, err
	}

	err = parseHeadersAndBody(lines[1:], &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func parseHeadersAndBody(lines []string, r *Request) error {
	bodyStarted := false
	body := ""

	for index, line := range lines {
		if line == "" { // headers and body are separated with /r/n/r/n - so the slice gives me ""
			bodyStarted = true
			continue
		}
		if bodyStarted {
			if index != len(lines)-1 {
				body += line + "\r\n"
			}
		} else {
			err := parseHTTPHeader(line, r)
			if err != nil {
				return err
			}
		}
	}
	r.Body = body
	return nil
}

func parseHTTPHeader(s string, r *Request) error {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}

	separatorIndex := strings.Index(s, ":")
	if separatorIndex == -1 {
		return errors.New("invalid http header")
	}

	key := strings.TrimSpace(s[0:separatorIndex])
	value := strings.TrimSpace(s[separatorIndex+1:])
	if key == "" {
		return errors.New("invalid http header")
	}

	keyLowered := strings.ToLower(key)

	r.Headers[keyLowered] = value
	return nil
}

func parseFirstLine(s string, r *Request) error {
	s = strings.TrimSpace(s)
	result := strings.Split(s, " ")

	len := len(result)
	if len != 3 {
		return errors.New("there must be 3 items in the first line")
	}

	method := result[0]
	err := parseMethod(method, r)
	if err != nil {
		return err
	}

	path := result[1]
	err = parsePath(path, r)
	if err != nil {
		return err
	}

	httpVersion := result[2]
	return parseHTTPVersion(httpVersion, r)
}

func parseMethod(s string, r *Request) error {
	httpMethod := HTTPMethod(s)
	switch httpMethod {
	case GET, POST, PUT, PATCH, DELETE:
		r.Method = httpMethod
	default:
		return errors.New("invalid http method")
	}
	return nil
}

func parsePath(s string, r *Request) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("path is empty")
	}

	splitPath := strings.SplitN(s, "?", 2)

	decodedPath, err := decodeHTTPPath(splitPath[0])
	if err != nil {
		return err
	}
	r.Path = decodedPath
	r.PathSegments = splitHTTPPathIntoSegments(decodedPath)

	if len(splitPath) > 1 {
		decodeQueryParams, err := decodeQueryParams(splitPath[1])
		if err != nil {
			return err
		}
		err = parseQueryParams(decodeQueryParams, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitHTTPPathIntoSegments(s string) []string {
	segments := strings.Split(s, "/")

	var result []string
	for _, segment := range segments {
		if segment != "" {
			result = append(result, segment)
		}
	}

	return result
}

func decodeQueryParams(s string) (string, error) {
	replaced := strings.ReplaceAll(s, "+", " ")
	return decodeHTTPPath(replaced)
}

func decodeHTTPPath(s string) (string, error) {
	if !strings.Contains(s, "%") {
		return s, nil
	}

	sLen := len(s)
	result := ""
	for i := 0; i < sLen; i++ {
		if s[i] == '%' {
			if (i + 2) >= sLen {
				return "", errors.New("invalid percent encoding")
			} else {
				encoded, err := hex.DecodeString(string(s[i+1]) + string(s[i+2]))
				if err != nil {
					return "", errors.New("invalid percent encoding")
				}
				result += string(encoded)
				i += 2
			}
		} else {
			result += string(s[i])
		}
	}

	return result, nil
}

func parseQueryParams(s string, r *Request) error {
	separatedParams := strings.SplitSeq(s, "&")
	for param := range separatedParams {
		if param == "" {
			continue
		}
		keyValue := strings.SplitN(param, "=", 2)
		if keyValue[0] == "" {
			return errors.New("query param key is empty")
		}
		if r.QueryParams == nil {
			r.QueryParams = make(map[string]string)
		}

		r.QueryParams[keyValue[0]] = ""
		if len(keyValue) > 1 {
			r.QueryParams[keyValue[0]] = keyValue[1]
		}
	}
	return nil
}

func parseHTTPVersion(s string, r *Request) error {
	versionLowered := strings.ToLower(s)
	if !strings.HasPrefix(versionLowered, "http/") {
		return errors.New("invalid http version")
	}

	versionNumber := versionLowered[5:]
	switch versionNumber {
	case "0.9", "1.0", "1.1", "2", "2.0", "3":
		r.Version = s
	default:
		return errors.New("invalid version number")
	}

	return nil
}
