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
	// listen on 127.0.0.1 on port server_port for client connections
	sock, err := net.Listen("tcp", "127.0.0.1:"+server_port)
	if err != nil {
		log.Fatalf("Could not setup a socket and listen on port: %s\n", err)
	}
	defer sock.Close()

	for {
		// accept any client connection
		conn, err := sock.Accept()
		if err != nil {
			log.Fatalf("Could not connect to client - %s\n", err)
		}
		// Close connection when this function ends
		defer conn.Close()

		// setup a reader to continuously read packets from the connection to a buffer so that we can print immediately
		// reader w/ buffer of size RECV_BUFFER_SIZE
		inp_stream := bufio.NewReaderSize(conn, RECV_BUFFER_SIZE)
		buffer := make([]byte, RECV_BUFFER_SIZE)

		// for loop handles when entire data is too long to be sent in one go, so it's sent in chunks
		for {
			// reads bytes from connection (up to RECV_BUFFER_SIZE) and stores them in buffer
			num_rec_bytes, err := inp_stream.Read(buffer)

			// EOF tells us we read everything
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving client data: %s\n", err)
			}
			// print all the bytes we received as soon as we receive them to stdout
			fmt.Print(string(buffer[:num_rec_bytes]))
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
