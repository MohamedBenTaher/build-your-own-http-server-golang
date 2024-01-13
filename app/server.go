package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var responseMap = map[string]string{
	"/": "HTTP/1.1 200 OK\r\n\r\n",
}

func handler(conn net.Conn) {
	request := make([]byte, 1024)
	_, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	requestData := strings.Split(string(request), " \r\n")
	path := strings.Split(requestData[0], " ")[1]
	if path, ok := responseMap[path]; ok {
		response := path
		conn.Write([]byte(response))
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Something went wrong while sending response")
		}
	} else {
		response := "HTTP/1.1 404 Not Found\r\n\r\n"
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Something went wrong while sending response")
			os.Exit(1)
		}
	}
}
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	handler(conn)
}
