package main

import (
	"fmt"
	"log"
	"net"

	"github.com/eskog/tcptohttp/internal/request"
)

func main() {
	listner, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatalf("Could not open file: %s", err)
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("error accepting a connection: %s", err)
			continue
		}
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("Error handling request: %s", err)
			continue
		}
		fmt.Printf(`
		Request line:
		- Method: %s
		- Target: %s
		- Version: %s
		`, req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
	}
}
