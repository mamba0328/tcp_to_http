package main

import (
	"fmt"
	"net"
	"tcp_to_http/internal/request"
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
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("request.RequestFromReader error: %v", err)
			continue
		}

		fmt.Printf("\tRequest line: \n\t\t- Method: %v\n\t\t- Target: %v\n\t\t- Version: %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		fmt.Println("Connection has been closed")
	}
}
