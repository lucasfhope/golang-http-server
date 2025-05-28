# Go HTTP Server

This is my implementation of a simple HTTP server, built in Go, by following the [Codecrafters "Build Your Own HTTP Server" challenge](https://app.codecrafters.io/courses/http-server/overview). The HTTP server supports both `GET` and `POST` methods, can handle multiple concurrent connections using goroutines, and supports persistent HTTP/1.1 connections, allowing multiple requests to be handled on the same connection. It also includes features like gzip encoding and file handling.

---

## Table of Contents

1. [Getting Started](#getting-started)
   - [Requirements](#requirements)
   - [Quickstart](#quickstart)
2. [Usage](#usage)



2. [Features](#features)  
3. [Endpoints](#endpoints)  
   - [GET Methods](#get-methods)  
   - [POST Methods](#post-methods)  
4. [Concurrency in Go](#concurrency-in-go)  
5. [How to Run](#how-to-run)  
6. [Testing the Server](#testing-the-server)  
7. [Future Improvements](#future-improvements)  

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

This server listens on port `4221` and processes HTTP requests using the TCP protocol, which ensures that data is delivered in the correct order without losing any data.

## Start the server

Run `app/main.go` to start the server.

```bash
go run app/main.go
```

## Interact with the server

You can interact with the server




---

## Features

- **Persistent HTTP/1.1 Connections**: Supports multiple requests on the same connection unless explicitly closed by the client.
- **Concurrency**: Handles multiple clients simultaneously using goroutines.
- **Gzip Encoding**: Automatically compresses responses if the client supports gzip encoding.
- **File Handling**: Allows uploading and downloading files via specific endpoints.
- **Custom Headers**: Processes headers like `User-Agent` and `Accept-Encoding`.

---

## Endpoints

### GET Methods

| **Target**            | **Description**                                                                 |
|------------------------|---------------------------------------------------------------------------------|
| `/`                   | Returns a `200 OK` response with no content.                                    |
| `/echo/<message>`      | Echoes back the `<message>` provided in the URL. Supports gzip encoding.        |
| `/user-agent`         | Returns the `User-Agent` header sent by the client.                             |
| `/files/<filename>`   | Serves the file specified by `<filename>` from the `./files` directory.          |

#### Example Requests:
1. **Echo a Message**:
   ```bash
   curl --http1.1 -v http://localhost:4221/echo/hello
