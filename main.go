package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLines(result chan string, file io.ReadCloser)  {
	defer close(result)
	defer file.Close()
	str := ""
	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
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
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", "error",  err)
	}
	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
