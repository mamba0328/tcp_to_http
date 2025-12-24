package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	// Open file
	listener, err := net.Listen("tcp", ":42069")
	defer listener.Close()

	if err != nil {
		fmt.Printf("net.Listen error: %v", err)
		return
	}

	for true {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("net.Listen error: %v", err)
			continue
		}

		fmt.Println("Connection has been accepted")

		// Get lines async
		lineChan := getLinesChannel(conn)

		for line := range lineChan {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("Connection has been closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lineChan := make(chan string)

	go func() {
		defer f.Close()
		defer close(lineChan)

		chunk := make([]byte, 8)
		line := ""

		for n, err := f.Read(chunk); ; n, err = f.Read(chunk) {
			if err != nil {
				if errors.Is(err, io.EOF) {
					lineChan <- line
					break
				}

				fmt.Printf("error: %s\n", err.Error())
				break
			}

			stringFromBytes := string(chunk[:n])
			stringParts := strings.Split(stringFromBytes, "\n")

			for n, part := range stringParts {
				if n == 0 {
					line += part
					continue
				}

				//Send line
				lineChan <- line
				// Reset prev line and store new line part
				line = part
			}
		}
	}()

	return lineChan
}
