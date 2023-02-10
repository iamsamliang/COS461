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
	// Connect to the server w/ server_ip at server_port
	conn, err := net.Dial("tcp", server_ip+":"+server_port)
	if err != nil {
		log.Fatalf("Could not connect to the server: %s\n", err)
	}

	// create a buffer to store the data from stdin so we can immediately send it to server
	// create a reader to read from stdin and put the data in buffer. Read at most SEND_BUFFER_SIZE
	buffer := make([]byte, SEND_BUFFER_SIZE)
	inp_stream := bufio.NewReaderSize(os.Stdin, SEND_BUFFER_SIZE)

	// for loop handles when data from stdin cannot be read in one go bc it is too big, so we need to send it in chunks
	for {
		// read from stdin and put the data in buffer
		bytes_read, err := inp_stream.Read(buffer)

		// EOF tells us we read everything from stdin, so we are done
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error reading from stdin: %s\n", err)
		}

		// send the data read into the buffer to the server
		conn.Write(buffer[:bytes_read])

		if err != nil {
			log.Fatalf("Error sending data to server: %s\n", err)
		}
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
