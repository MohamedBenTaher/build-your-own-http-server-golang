# Web Server Documentation

## Introduction
This is a simple web server written in Go that serves files from a specified directory.

The web supports is inspired by the [Codecrafters HTTP Server Challenge](https://codecrafters.io/challenges/http-server) for learning Golang and understanding how web servers work under the hood.

## Usage
To use the web server, follow these steps:

1. Build the executable by running the following command:
    ```
    go build
    ```

2. Run the executable with the following command:
    ```
    ./webserver -directory <directory_path>
    ```
    Replace `<directory_path>` with the path to the directory you want to serve. If no directory is specified, the current directory will be served.

3. The web server will start listening on port 4221. You can access it by opening a web browser and navigating to `http://localhost:4221`.

## API Endpoints

### GET /user-agent
Returns the User-Agent header of the HTTP request.

#### Request
- Method: GET
- Path: /user-agent

#### Response
- Status Code: 200 OK
- Content-Type: text/plain
- Body: The User-Agent header value

### GET /echo/{message}
Returns the provided message as the response body.

#### Request
- Method: GET
- Path: /echo/{message}
- Replace `{message}` with the desired message.

#### Response
- Status Code: 200 OK
- Content-Type: text/plain
- Body: The provided message

### GET /files/{file_name}
Returns the contents of the specified file.

#### Request
- Method: GET
- Path: /files/{file_name}
- Replace `{file_name}` with the name of the file you want to retrieve.

#### Response
- Status Code: 200 OK
- Content-Type: application/octet-stream
- Body: The contents of the file

### POST /files/{file_name}
Creates a new file with the provided name and content.

#### Request
- Method: POST
- Path: /files/{file_name}
- Replace `{file_name}` with the name of the file you want to create.
- Body: The content of the file

#### Response
- Status Code: 201 Created
- Content-Type: application/octet-stream
- Body: "File created successfully"

## Error Handling
The web server handles the following errors:

- 404 Not Found: If the requested resource is not found.
- 500 Internal Server Error: If there is an internal server error.

## Dependencies
The web server uses the following dependencies:

- `net` package: For network operations.
- `flag` package: For command-line flag parsing.
- `os` package: For file operations.
- `path/filepath` package: For working with file paths.
- `strings` package: For string manipulation.
- `bufio` package: For reading from network connections.

## License
This web server is released under the MIT License. See the [LICENSE](LICENSE) file for more details.
