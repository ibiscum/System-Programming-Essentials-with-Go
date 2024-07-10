package main

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	// The native way to print "Hello, World!" to stdout
	fmt.Println("Hello, World!")

	// The overly complicated way to print "Hello, World!" to stdout
	_, _, err := unix.Syscall(unix.SYS_WRITE, 1,
		uintptr(unsafe.Pointer(&[]byte("Hello, World!")[0])),
		uintptr(len("Hello, World!")),
	)
	if err != 0 {
		log.Fatal(err)
	}
}
