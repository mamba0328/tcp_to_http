package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
