package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       int // 0 Initialized, 1 Done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	fullRequest, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, _, err := parseRequestLine(fullRequest)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(request []byte) (RequestLine, int, error) {
	requestString := string(request)
	requestLineString := strings.Split(requestString, "\r\n")[0]

	if requestLineString == "" {
		return RequestLine{}, 0, nil
	}

	requestLineParts := strings.Split(requestLineString, " ")

	if len(requestLineParts) != 3 {
		return RequestLine{}, len(request), errors.New("Bad request line")
	}

	httpVersionParts := strings.Split(requestLineParts[2], "/")

	if len(httpVersionParts) != 2 {
		return RequestLine{}, len(request), errors.New("Bad http version")
	}

	if httpVersionParts[0] != "HTTP" {
		return RequestLine{}, len(request), fmt.Errorf("unrecognized HTTP-version: %s", httpVersionParts[0])
	}

	requestLine := RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersionParts[1],
	}

	err := validateRequestLine(requestLine)

	if err != nil {
		return RequestLine{}, len(request), err
	}

	return requestLine, len(request), nil
}

func validateRequestLine(rl RequestLine) error {
	validMethods := []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
	}

	if !slices.Contains(validMethods, rl.Method) {
		return errors.New("Invalid Request Method")
	}

	if rl.HttpVersion != "1.1" {
		return errors.New("Unsupported Http Version")
	}

	// Target validation

	return nil
}
