package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
)

func filterLogs(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logEntry := scanner.Text()
		if strings.Contains(logEntry, "ERROR") {
			_, err := writer.Write([]byte(logEntry + "\n"))
			checkError(err)
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	pipePath := "/tmp/my_log_pipe"
	if err := os.RemoveAll(pipePath); err != nil {
		panic(err)
	}
	if err := syscall.Mkfifo(pipePath, 0600); err != nil {
		panic(err)
	}
	defer os.RemoveAll(pipePath)

	pipeFile, err := os.OpenFile(pipePath, os.O_RDONLY|os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	defer pipeFile.Close()

	go func() {
		writer, err := os.OpenFile(pipePath, os.O_WRONLY, os.ModeNamedPipe)
		checkError(err)
		defer writer.Close()

		for {
			_, err = writer.WriteString("INFO: All systems operational\n")
			checkError(err)

			_, err = writer.WriteString("ERROR: An error occurred\n")
			checkError(err)

			time.Sleep(1 * time.Second)
		}
	}()

	filterLogs(pipeFile, os.Stdout)
}
