package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	defaultTarget   = "/"
	echoTarget      = "/echo/"
	userAgentTarget = "/user-agent"
	filesTarget     = "/files/"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	fmt.Println("Listening on port 4221...")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Print("Error accepting connection:\n", err.Error(), "\n")
			continue
		}
		go handleConnection(conn)
	}
}

///////////////////////////////////
// Connection Request Handling  //
//////////////////////////////////

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	for {
		reader := bufio.NewReader(conn)
		request, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Println("Error reading request:", err.Error())
			continue
		}

		requestArray := strings.Split(request, " ")
		if len(requestArray) != 3 {
			fmt.Println("Invalid request format")
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"))
			continue
		}
		method := requestArray[0]
		target := requestArray[1]
		fmt.Printf("Received %s request for target %s\n", method, target)

		headers := make(map[string]string)
		var errorReadingHeaders bool
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				errorReadingHeaders = true
				fmt.Print("Error reading header\n", err.Error(), "\n")
				continue
			}
			line = strings.TrimSpace(line)

			if line == "" {
				break
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				headers[key] = value
			}
		}
		if errorReadingHeaders {
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"))
			continue
		}

		body := []byte{}
		contentLengthStr, hasBody := headers["Content-Length"]
		if hasBody {
			var contentLength int
			_, err := fmt.Sscanf(contentLengthStr, "%d", &contentLength)
			if err != nil {
				fmt.Println("Invalid Content-Length:", contentLengthStr)
				conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"))
				continue
			}

			body = make([]byte, contentLength)
			totalRead := 0
			for totalRead < contentLength {
				n, err := reader.Read(body[totalRead:])
				if err != nil {
					fmt.Println("Error reading body:", err)
					conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 0\r\n\r\n"))
					continue
				}
				totalRead += n
			}
		}

		if method == "GET" {
			handleGetRequest(conn, target, headers, body)
		} else if method == "POST" {
			handlePostRequest(conn, target, headers, body)
		} else {
			conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\nAllow: GET, POST\r\n\r\n"))
		}

		connectionHeader, exists := headers["Connection"]
		if exists && strings.ToLower(connectionHeader) == "close" {
			fmt.Println("Connection: close received, closing connection")
			return
		}
	}
}

func handleGetRequest(conn net.Conn, target string, headers map[string]string, body []byte) {
	if target == defaultTarget {
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
	} else if strings.HasPrefix(target, echoTarget) {
		getEcho(conn, target, headers)
	} else if target == userAgentTarget {
		getUserAgent(conn, headers)
	} else if strings.HasPrefix(target, filesTarget) {
		getFile(conn, target)
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))
	}
}

func handlePostRequest(conn net.Conn, target string, headers map[string]string, body []byte) {
	if strings.HasPrefix(target, filesTarget) {
		postFile(conn, target, body)
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))
	}
}

////////////////////////////
// GET Request Functions  //
////////////////////////////

func getEcho(conn net.Conn, target string, headers map[string]string) {
	echoMessage := strings.TrimSpace(strings.TrimPrefix(target, echoTarget))
	contentEncoding, encodedData, err := getEncodingData(headers, echoMessage)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		return
	}
	var responseBody []byte
	if encodedData != nil {
		responseBody = encodedData
	} else {
		responseBody = []byte(echoMessage)
	}
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n%sContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n",
		contentEncoding,
		len(responseBody),
	)
	conn.Write([]byte(response))
	conn.Write(responseBody)
}

func getUserAgent(conn net.Conn, headers map[string]string) {
	userAgent, exists := headers["User-Agent"]
	if !exists {
		fmt.Println("User-Agent header not found")
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))
		return
	}
	userAgent = strings.TrimSpace(userAgent)
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
	conn.Write([]byte(response))
}

func getFile(conn net.Conn, target string) {
	filePath := "." + target
	fmt.Println("File path:", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Print("Error reading file:\n", err.Error(), "\n")
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"))
		return
	}
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
	conn.Write([]byte(response))
}

////////////////////////////
// POST Request Handlers  //
////////////////////////////

func postFile(conn net.Conn, target string, body []byte) {
	fileName := strings.TrimSpace(strings.TrimPrefix(target, filesTarget))

	// Create the files directory if it doesn't exist
	dirPath := "./files"
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		fmt.Print("Error creating directory\n", err.Error(), "\n")
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\nContent-Length: 0\r\n\r\n"))
		return
	}

	// Create the file
	filePath := dirPath + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Print("Error creating file:\n", err.Error(), "\n")
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\nContent-Length: 0\r\n\r\n"))
		return
	}
	defer file.Close()

	fmt.Println(body)

	// Write the body to the file
	_, err = file.Write(body)
	if err != nil {
		fmt.Print("Error writing to file\n", err.Error(), "\n")
		conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\nContent-Length: 0\r\n\r\n"))
		return
	}

	fmt.Println("File created successfully:", filePath)
	conn.Write([]byte("HTTP/1.1 201 Created\r\nContent-Length: 0\r\n\r\n"))
}

//////////////////////
// gzip Encoding  ////
//////////////////////

func getEncodingData(headers map[string]string, text string) (string, []byte, error) {
	value, exists := headers["Accept-Encoding"]
	if !exists {
		return "", nil, nil
	}
	for _, encoding := range strings.Split(value, ",") {
		if strings.Contains(strings.TrimSpace(encoding), "gzip") {
			encodedText, err := gzipEncode(text)
			if err != nil {
				fmt.Print("Error gzip encoding text\n", err.Error(), "\n")
				return "", nil, err
			}
			return "Content-Encoding: gzip\r\n", encodedText, nil
		}
	}
	return "", nil, nil
}

func gzipEncode(data string) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	defer writer.Close()

	_, err := writer.Write([]byte(data))
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
