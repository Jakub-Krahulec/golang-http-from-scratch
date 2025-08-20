package main

import (
	"encoding/hex"
	"errors"
	"strings"
)

func parseRequest(s string) (*Request, error) {
	lines := strings.Split(s, "\r\n")
	if len(lines) == 0 {
		return nil, errors.New("there arent any lines") // TODO: make better error message later
	}

	var r Request
	err := parseFirstLine(lines[0], &r)
	if err != nil {
		return nil, err
	}

	loopingOverHeaders := true
	body := ""
	for index, line := range lines {
		if index == 0 {
			continue
		}
		if line == "" { // hledam dvakrat /r/n za sebou... ale slice mi to rozdelil a udelal z toho teda empty string
			loopingOverHeaders = false
		}
		if loopingOverHeaders {
			err := parseHTTPHeader(line, &r)
			if err != nil {
				return nil, err
			}
		} else {
			body += line + "\r\n" // pridavam nakonec zbytecne novej radek
		}
	}
	r.Body = body
	return &r, nil
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
		return errors.New("there must be 3 items in first line") // TODO: make better error message when figure out whole structure
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
