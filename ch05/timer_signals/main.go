package main

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/unix" // This package is used to interact with the Linux kernel, if you're on a different OS, you'll need to use a different package
)

func clen(n []byte) int {
	for i, b := range n {
		if b == 0 {
			return i
		}
	}
	return len(n)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fd, err := unix.InotifyInit()
	if err != nil {
		fmt.Println("Error initializing inotify:", err)
		return
	}
	defer unix.Close(fd)

	watchPath := "/path/to/your/directory" // Change this to the directory you want to watch
	watchDescriptor, err := unix.InotifyAddWatch(fd, watchPath, unix.IN_MODIFY|unix.IN_CREATE|unix.IN_DELETE)
	if err != nil {
		fmt.Println("Error adding watch:", err)
		return
	}
	defer func() {
		_, err := unix.InotifyRmWatch(fd, uint32(watchDescriptor))
		checkError(err)
	}()

	const bufferSize = (unix.SizeofInotifyEvent + unix.NAME_MAX + 1)
	buf := make([]byte, bufferSize)
	for {
		n, err := unix.Read(fd, buf[:])
		if err != nil {
			fmt.Println("Error reading from inotify:", err)
			return
		}

		var offset uint32
		for offset < uint32(n) {
			event := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))
			nameBytes := buf[offset+unix.SizeofInotifyEvent : offset+unix.SizeofInotifyEvent+uint32(event.Len)]
			name := string(nameBytes[:clen(nameBytes)])

			fmt.Printf("Event: %s/%s\n", watchPath, name)
			offset += unix.SizeofInotifyEvent + uint32(event.Len)
		}
	}
}
