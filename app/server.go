package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	// Uncomment this block to pass the first stage
	dirFlag := flag.String("directory", ".", "The directory to serve")
	flag.Parse()
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
		go handleConn(conn, dirFlag)
	}
}

type Request struct {
	HttpRequest
	RequestHeaders
	RequestBody string
	Accept      string
}
type HttpRequest struct {
	Method  string
	Target  string
	Version string
}
type RequestHeaders struct {
	Host      string
	UserAgent string
}
type ResponseHeaders struct {
	ContentType   string
	ContentLength int
}
type ResponseBody struct {
	Message string
}
type Response struct {
	Version       string
	StatusCode    int
	StatusMessage string
	ResponseHeaders
	ResponseBody
}

func NewRequest(conn net.Conn) Request {
	var httpReq HttpRequest
	var reqHeaders RequestHeaders
	var requestBody string
	l := bufio.NewReader(conn)
	for {
		line, err := l.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			panic(err)
		}
		if line == "\r\n" {
			requestBody, _ = l.ReadString('\n')
			break
		}
		if strings.HasPrefix(line, "User-Agent:") {
			fmt.Println(line)
			s := strings.Split(line, " ")[1]
			reqHeaders.UserAgent = s
		}
		if strings.HasPrefix(line, "Host:") {
			fmt.Println(line)
			s := strings.Split(line, " ")[1]
			reqHeaders.Host = s
		}
		if strings.Contains(line, "HTTP/1.1") {
			fmt.Println(line)
			s1 := strings.Split(line, " ")[0]
			s2 := strings.Split(line, " ")[1]
			s3 := strings.Split(line, " ")[2]
			httpReq.Version = strings.TrimSuffix(s3, "\r\n")
			httpReq.Method = s1
			httpReq.Target = s2
		}
	}
	return Request{HttpRequest: httpReq, RequestHeaders: reqHeaders, RequestBody: requestBody, Accept: "*/*"}
}
func (r *Request) print() {
	fmt.Println(r)
}
func (r *Request) parseUserAgent() Response {
	if r.Method == "GET" && r.Target == "/user-agent" {
		fmt.Println(r.UserAgent)
		trimUserAgent := strings.TrimSpace(r.UserAgent)
		return Response{
			Version:       r.Version,
			StatusCode:    200,
			StatusMessage: "OK",
			ResponseHeaders: ResponseHeaders{
				ContentType:   "text/plain",
				ContentLength: len(trimUserAgent),
			},
			ResponseBody: ResponseBody{
				Message: trimUserAgent,
			},
		}
	} else {
		return Response{
			Version:       r.Version,
			StatusCode:    404,
			StatusMessage: "Not Found",
			ResponseHeaders: ResponseHeaders{
				ContentType:   "text/plain",
				ContentLength: len(r.UserAgent),
			},
			ResponseBody: ResponseBody{
				Message: r.UserAgent,
			},
		}
	}
}
func (r *Request) handleEcho() Response {
	respStr := strings.TrimPrefix(r.Target, "/echo/")
	fmt.Println(respStr)
	return Response{
		Version:       r.Version,
		StatusCode:    200,
		StatusMessage: "OK",
		ResponseHeaders: ResponseHeaders{
			ContentType:   "text/plain",
			ContentLength: len(respStr),
		},
		ResponseBody: ResponseBody{
			Message: respStr,
		},
	}
}

func (r *Request) checkFiles(fileWithPath string) bool {
	wd, _ := os.Getwd()
	fmt.Println("Current working directory:", wd)
	fmt.Println("Checking file:", fileWithPath)

	fileStat, err := os.Stat(fileWithPath)
	print(fileStat)
	if err != nil {
		fmt.Println("Error accessing file: ", err.Error())
		return false
	}
	if fileStat.IsDir() {
		fmt.Println("File name pass is a directory")
		return false
	} else {
		return true
	}
}
func (r *Request) return404() Response {
	return Response{
		Version:       r.Version,
		StatusCode:    404,
		StatusMessage: "Not Found",
		ResponseHeaders: ResponseHeaders{
			ContentType:   "text/plain",
			ContentLength: 0,
		},
		ResponseBody: ResponseBody{
			Message: "",
		},
	}
}
func (r *Request) return500() Response {
	return Response{
		Version:       r.Version,
		StatusCode:    500,
		StatusMessage: "Internal Server Error",
		ResponseHeaders: ResponseHeaders{
			ContentType:   "text/plain",
			ContentLength: 0,
		},
		ResponseBody: ResponseBody{
			Message: "",
		},
	}
}
func (r *Request) handleFGetiles(dirFlag *string) Response {
	fileName := strings.TrimPrefix(r.Target, "/files/")
	filePlusPath := filepath.Join(*dirFlag, fileName)
	if r.checkFiles(filePlusPath) {
		// do stuff with files
		fileContents, err := os.ReadFile(filePlusPath)
		if err != nil {
			fmt.Println("Error reading file: ", err.Error())
			// return 404
			return r.return404()
		}
		fmt.Println("File contents:", fileContents)
		return Response{
			Version:       r.Version,
			StatusCode:    200,
			StatusMessage: "OK",
			ResponseHeaders: ResponseHeaders{
				ContentType:   "application/octet-stream",
				ContentLength: len(fileContents),
			},
			ResponseBody: ResponseBody{
				Message: string(fileContents),
			},
		}
	} else {
		// return 404
		return r.return404()
	}
}
func (r *Request) handleFPostiles(dirFlag *string) Response {
	fileName := strings.TrimPrefix(r.Target, "/files/")
	filePlusPath := filepath.Join(*dirFlag, fileName)

	// Create the file
	file, err := os.Create(filePlusPath)
	if err != nil {
		fmt.Println("Error creating file: ", err.Error())
		// return 500
		return r.return500()
	}
	defer file.Close()

	// Write the request body to the file
	_, err = file.WriteString(r.RequestBody)
	if err != nil {
		fmt.Println("Error writing to file: ", err.Error())
		// return 500
		return r.return500()
	}

	// return 201 Created
	return Response{
		Version:       r.Version,
		StatusCode:    201,
		StatusMessage: "Created",
		ResponseHeaders: ResponseHeaders{
			ContentType:   "application/octet-stream",
			ContentLength: len(r.RequestBody),
		},
		ResponseBody: ResponseBody{
			Message: "File created successfully",
		},
	}
}
func (r *Request) WriteResponse(resp Response, conn net.Conn) {
	fmt.Println(resp)
	newResp := fmt.Sprintf(
		"%s %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		resp.Version,
		resp.StatusCode,
		resp.StatusMessage,
		resp.ResponseHeaders.ContentType,
		resp.ResponseHeaders.ContentLength,
		resp.ResponseBody.Message,
	)
	fmt.Println(newResp)
	conn.Write([]byte(newResp))
}
func (r *Request) RouteRequest(dirFLag *string) Response {
	switch {
	case r.Target == "/user-agent":
		return r.parseUserAgent()
	case strings.Contains(r.Target, "/echo/"):
		return r.handleEcho()
	case r.Target == "/":
		return Response{
			Version:       r.Version,
			StatusCode:    200,
			StatusMessage: "OK",
			ResponseHeaders: ResponseHeaders{
				ContentType:   "text/plain",
				ContentLength: 0,
			},
			ResponseBody: ResponseBody{
				Message: "",
			},
		}
	case strings.Contains(r.Target, "/files/"):
		if r.Method == "GET" {
			return r.handleFGetiles(dirFLag)
		} else if r.Method == "POST" {
			return r.handleFPostiles(dirFLag)
		}
	}
	// Default return statement
	return Response{
		Version:       r.Version,
		StatusCode:    404,
		StatusMessage: "Not Found",
		ResponseHeaders: ResponseHeaders{
			ContentType:   "text/plain",
			ContentLength: 0,
		},
		ResponseBody: ResponseBody{
			Message: "",
		},
	}
}
func handleConn(conn net.Conn, dirFlag *string) {
	newR := NewRequest(conn)
	// newR.print()
	// parseAgent := newR.parseUserAgent()
	resp := newR.RouteRequest(dirFlag)
	newR.WriteResponse(resp, conn)
}
