package main

import (
	"bufio"
	"bytes"
	"log"
)

func main() {
	buf := bufio.NewReader(bytes.NewBuffer([]byte{0x59, 0x44}))
	p := make([]byte, 0, 0)
	n, err := buf.Read(p)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Read: ", p, "Count:", n)

}
