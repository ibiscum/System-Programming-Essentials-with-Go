package main

import (
	"io"
	"log"
	"net"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		byteP := make([]byte, 1024)
		return &byteP
	},
}

func handleConnection(conn net.Conn) {
	buf := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(&buf)

	for {
		n, err := conn.Read(*buf)
		if err != nil {
			if err != io.EOF {
				println("Error reading:", err.Error())
			}
			break
		}

		tmp_buf := *buf
		_, err = conn.Write(tmp_buf[:n])
		if err != nil {
			log.Fatal(err)
		}
	}
	conn.Close()
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accepting connection:", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}
