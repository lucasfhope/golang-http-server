# Golang HTTP Server

This is a simple HTTP server implemented in Go that supports both `GET` and `POST` methods. The server handles multiple concurrent connections using Go's goroutines and supports persistent HTTP/1.1 connections. It also includes features like gzip encoding and file handling.

---

## Table of Contents

1. [Overview](#overview)  
2. [Features](#features)  
3. [Endpoints](#endpoints)  
   - [GET Methods](#get-methods)  
   - [POST Methods](#post-methods)  
4. [Concurrency in Go](#concurrency-in-go)  
5. [How to Run](#how-to-run)  
6. [Testing the Server](#testing-the-server)  
7. [Future Improvements](#future-improvements)  

---

## Overview

This server listens on port `4221` and processes HTTP requests. It supports persistent connections, allowing multiple requests to be handled on the same connection. The server is designed to handle concurrent connections efficiently using Go's concurrency model.

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
