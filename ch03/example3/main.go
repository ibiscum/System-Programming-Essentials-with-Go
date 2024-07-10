package main

import (
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	_, err := unix.Write(1, []byte("Hello, World!"))
	if err != nil {
		log.Fatal(err)
	}
}
