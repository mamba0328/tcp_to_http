package request

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	State       requestState // 0 Initialized, 1 Done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const bufferSize = 8

func (r *Request) parse(data []byte) (int, error) {
	if r.State == requestStateDone {
		return 0, errors.New("parse error: trying to parse during Request `Done` state")
	}

	if r.State != 0 {
		return 0, errors.New("parse error: unknown Request state")
	}

	requestLine, n, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}
	r.RequestLine = requestLine
	r.State = requestStateDone

	return n, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0
	buffer := make([]byte, bufferSize)

	request := &Request{
		State: requestStateInitialized,
	}

	for request.State == requestStateInitialized {
		if len(buffer) >= readToIndex {
			extendedBuffer := make([]byte, len(buffer)*2)
			copy(extendedBuffer, buffer)
			buffer = extendedBuffer
		}
		readCount, err := reader.Read(buffer[readToIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = requestStateDone
				break
			}
			return nil, err
		}

		readToIndex += readCount

		parseCount, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		// HACK - CLEARS BUFFER IF WE HAVE PARSE COUNT
		copy(buffer, buffer[parseCount:])
		readToIndex -= parseCount
	}

	return request, nil
}

func parseRequestLine(request []byte) (RequestLine, int, error) {
	//WAIT FOR FULL REQUEST STRING
	requestString := string(request)

	requestParts := strings.Split(requestString, "\r\n")

	if len(requestParts) <= 1 {
		return RequestLine{}, 0, nil
	}

	requestLineString := requestParts[0]

	requestLine, err := parseRequestLineFromString(requestLineString)

	if err != nil {
		return RequestLine{}, 0, err
	}

	return requestLine, len(requestLineString) + 2, nil
}

func parseRequestLineFromString(requestLineString string) (RequestLine, error) {
	fmt.Println(requestLineString)
	requestLineParts := strings.Split(requestLineString, " ")

	if len(requestLineParts) != 3 {
		return RequestLine{}, errors.New("not enough request line parts")
	}

	httpVersionParts := strings.Split(requestLineParts[2], "/")

	if len(httpVersionParts) != 2 {
		return RequestLine{}, errors.New("not enough request http version parts")
	}

	if httpVersionParts[0] != "HTTP" {
		return RequestLine{}, fmt.Errorf("unrecognized HTTP-version: %s", httpVersionParts[0])
	}

	requestLine := RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersionParts[1],
	}

	err := validateRequestLine(requestLine)

	if err != nil {
		return RequestLine{}, err
	}

	return requestLine, nil
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
