package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"medivhtcp/internal/request"
	"net"
)

func getLines(result chan string, r io.ReadCloser) {
	defer close(result)
	defer r.Close()
	str := ""
	for {
		data := make([]byte, 8)
		n, err := r.Read(data)
		if err != nil {
			break
		}
		data = data[:n]
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			str += string(data[:i])
			data = data[i+1:]

			result <- str
			str = ""
		}
		str += string(data)
	}
	if len(str) != 0 {
		result <- str
	}
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}
