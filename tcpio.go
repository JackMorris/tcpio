// tcpio - simple program to connect to a remote tcp socket and send/receieve data between
// the remote port and stdin/out.
// Exits on any failure (eg. connection closed, either side closed, ...)
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// Get the remote host.
	if len(os.Args) != 2 {
		fmt.Println("Usage: tcpio <ip>:<port>")
		os.Exit(1)
	}

	// Connect to the remote host.
	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Copy data both ways.
	done := make(chan struct{})
	go passData(os.Stdout, conn, done) // conn -> stdout
	go passData(conn, os.Stdin, done)  // stdin -> conn

	// Exit when either data pass is done.
	<-done
	conn.Close()
}

func passData(dst io.Writer, src io.Reader, done chan<- struct{}) {
	// Pass data from src to dst, and print log on error.
	// Signal on `done` when complete (whether or not from error).
	if _, err := io.Copy(dst, src); err != nil {
		log.Print(err)
	}
	done <- struct{}{}
}
