package main

import (
	"bytes"
	"log"
)

func main() {
	data := []byte("this is some data stored as a byte slice in Go Lang!")

	// convert byte slice to io.Reader
	reader := bytes.NewReader(data)

	// read only 4 byte from our io.Reader
	buf := make([]byte, 100000)
	n, err := reader.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(buf[:n]))
}
