/*****************************************************************************
 * client-go.go
 * Name: Sam Liang
 * NetId: saml
 *****************************************************************************/

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

const SEND_BUFFER_SIZE = 2048

/* TODO: client()
 * Open socket and send message from stdin.
 */
func client(server_ip string, server_port string) {
	conn, err := net.Dial("tcp", server_ip+":"+server_port)
	if err != nil {
		log.Fatalf("failed to connect to server - %v\n", err)
	}

	// inspired by this source code: https://gist.github.com/rodkranz/90c82583987a15e3d0f2c4678f2835c7
	buf := make([]byte, SEND_BUFFER_SIZE)
	reader := bufio.NewReaderSize(os.Stdin, SEND_BUFFER_SIZE)

	// for loop allows us to read multiple times (in case of big inputs)
	for {
		num_bytes, err := reader.Read(buf)

		// client should break when we hit EOF
		if err == io.EOF {
			break
		}

		// continue to process buf
		if err != nil {
			log.Fatalf("Failed to read %v\n", err)
		}

		conn.Write(buf[:num_bytes])
	}
}

// Main parses command-line arguments and calls client function
func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: ./client-go [server IP] [server port] < [message file]")
	}
	server_ip := os.Args[1]
	server_port := os.Args[2]
	client(server_ip, server_port)
}
