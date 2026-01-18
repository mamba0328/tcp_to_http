package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headerString := strings.Trim(string(data), " ")

	// Check is the end of the headers
	if headerString == "\r\n" {
		return 0, true, nil
	}

	// Wait for full line
	if !strings.Contains(headerString, "\r\n") {
		return 0, false, nil
	}

	// Check either Parse is done after this data
	//hasClosing := strings.Contains(headerString, "\r\n\r\n")

	fieldLine := strings.Split(headerString, "\r\n")[0]
	fieldLineSeparatorIndex := strings.Index(fieldLine, ":")

	fieldName := fieldLine[:fieldLineSeparatorIndex]
	fieldValue := fieldLine[fieldLineSeparatorIndex+1:]

	if strings.HasSuffix(fieldName, " ") {
		return 0, false, errors.New("headers.Parse: invalid field name")
	}

	h[fieldName] = strings.Trim(fieldValue, " ")

	return len(data), false, nil
}

func NewHeaders() Headers {
	return Headers{}
}
