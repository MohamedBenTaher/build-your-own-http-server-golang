package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var responseMap = map[string]string{
	"/": "HTTP/1.1 200 OK\r\n\r\nHello, world!",
}

func handler(conn net.Conn) {
	request := make([]byte, 1024)
	_, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	requestData := strings.Split(string(request), "\r\n")
	requestLine := strings.Split(requestData[0], " ")
	path := requestLine[1]
	if response, ok := responseMap[path]; ok {
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Something went wrong while sending response")
		}
	} else if strings.Contains(path, "/echo/") {
		body := strings.Split(path, "/echo/")[1]
		length := len(body)
		headers := []string{
			"HTTP/1.1 200 OK",
			fmt.Sprintf("Content-Type: text/plain; charset=utf-8"),
			fmt.Sprintf("Content-Length: %v", length),
		}
		response := fmt.Sprintf("%s\r\n\r\n%s", strings.Join(headers, "\r\n"), body)
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Something went wrong while sending response Body")
			os.Exit(1)
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handler(conn)
	}
}
