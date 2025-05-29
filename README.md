# Go HTTP Server

This is my implementation of a simple HTTP server, built in Go, by following the [Codecrafters "Build Your Own HTTP Server" challenge](https://app.codecrafters.io/courses/http-server/overview). The HTTP server, which uses the TCP protocol for realiable data transfer, supports both `GET` and `POST` methods, can handle multiple concurrent connections using goroutines, and supports persistent HTTP/1.1 connections, allowing multiple requests to be handled on the same connection. It also includes features like gzip encoding and file handling.

---

## Table of Contents

1. [Getting Started](#getting-started)  
   - [Requirements](#requirements)  
   - [Quickstart](#quickstart)  
2. [Usage](#usage)  
   - [Start the Server](#start-the-server)  
   - [Interact with the Server](#interact-with-the-server)  
      - [GET Endpoints](#get-endpoints)  
      - [POST Endpoint](#post-endpoint)  
3. [Other Features](#other-features)  
   - [Gzip Compression](#gzip-compression)  
   - [Persistent Connection](#persistent-connection)

---

# Getting Started

## Requirements

- **git**
    - Try running `git --version` to see if it is installed
- **go**
    - Try running `go version` to se if it is installed

## Quickstart

```bash
git clone https://github.com/lucasfhope/go-http-server.git
cd go-http-server
```

---

# Usage

This server listens on port `4221`.

## Start the server

Run `app/main.go` to start the server.

```bash
go run app/main.go
```

## Interact with the server

In another terminal, you can interact with the server by sending requests to `localhost:4221`.

```bash
curl -i localhost:4221
```

You should receive:

```
HTTP/1.1 200 OK
Content-Length: 0
```

The server can handle both `Get` and `POST` requests.

### GET Endpoints

| **Target**            | **Description**                                                                 |
|------------------------|---------------------------------------------------------------------------------|
| `/`                   | Returns a `200 OK` response with no content.                                    |
| `/echo/<message>`      | Echoes back the `<message>` provided in the URL. Supports gzip encoding.        |
| `/user-agent`         | Returns the `User-Agent` header sent by the client.                             |
| `/files/<filename>`   | Serves the file specified by `<filename>` from the `./files` directory.          |


### POST Endpoint

| **Target**            | **Description**                                                                 |
|------------------------|---------------------------------------------------------------------------------|
| `/files/<filename>`    | Creates a file with the specified <filename> and with the contents of the POST request body in it. |

---

# Other Features

## Gzip Compression

The HTTP server supports gzip compression. If send a `GET` request to the `echo` endpoint and include `Accept-Encoding: gzip` in the request header, then the server will send a response body with the content gzip encoded. You can use the following command.

```bash
curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/hello | hexdump -C
```

## Persistent Connection

The server can handle persistent connections, meaning that the same TCP connection can be used for multiple requests. You can test this with the following command.

```bash
curl --http1.1 -v http://localhost:4221/user-agent -H "User-Agent: orange/mango-grape" --next http://localhost:4221/echo/apple
```

Because of goroutines, the HTTP server is able to handle multiple persistent connections concurrently.

