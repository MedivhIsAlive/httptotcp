package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func getLines(result chan string, r io.ReadCloser)  {
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
			data = data[i + 1:]

			result <- str
			str = ""
		}
		str += string(data)
	}
	if len(str) != 0 {
		result <- str
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)
	go getLines(out, f)
	return out
}


func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Listening on :42069")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Could not accept connection")
			}
			fmt.Println("Connection accepted")
			go func(c net.Conn) {
				defer conn.Close()
				for line := range getLinesChannel(conn) {
					fmt.Printf("%s\n", line)
				}
				fmt.Println("Connection stopped")
			}(conn)
		}
	}()
	<-sig
	fmt.Println("\nWe're going down down...")
	os.Exit(0)
}
