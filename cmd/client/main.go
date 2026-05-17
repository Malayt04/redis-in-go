package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/Malayt04/redis-in-go/internal/resp"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command: ")
		command, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}

		encoded := resp.Encode(command)
		fmt.Println("Sending:", encoded)

		_, err = conn.Write([]byte(encoded))
		if err != nil {
			panic(err)
		}

		respReader := bufio.NewReader(conn)
		response, err := resp.DecodeResponse(respReader)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		if response.Type == "error" {
			fmt.Println("ERROR:", response.Value)
		} else if response.Type == "null" {
			fmt.Println("(nil)")
		} else {
			fmt.Println("Received:", response.Value)
		}
	}
}