/*****************************************************************************
 * server-go.go
 * Name: Sam Liang
 * NetId: saml
 *****************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const RECV_BUFFER_SIZE = 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 */
func server(server_port string) {
	ln, err := net.Listen("tcp", "localhost"+":"+server_port)
	if err != nil {
		log.Fatalf("Failed to setup a listener - %v\n", err)
	}
	defer ln.Close()

	for {
		// accept connection
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection - %v\n", err)
		}
		// Close connection when this function ends
		defer conn.Close()

		reader := bufio.NewReaderSize(conn, RECV_BUFFER_SIZE)
		buffer := make([]byte, RECV_BUFFER_SIZE)

		for {
			num_bytes, err := reader.Read(buffer)

			// keep reading until we hit an EOF
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read %v\n", err)
			}
			fmt.Print(string(buffer[:num_bytes]))
		}
	}
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
