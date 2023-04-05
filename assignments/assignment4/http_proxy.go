/*****************************************************************************
 * http_proxy.go
 * Names: Sam Liang
 * NetIds: saml
 *****************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("In handler method")
	if req.Method != "GET" {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Absolute URL: ", req.URL.String())
	relative_url := req.URL.Path
	fmt.Println("relative url: ", relative_url)
	header := req.Header
	fmt.Println("Original request header: ", header)

	// this URL should be a relative URL and not the absolute URL
	new_req, err := http.NewRequest("GET", relative_url, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error creating request: ", err)
		return
	}

	// new_req.URL, _ = url.Parse(req.URL.Path)

	// Set the header values in the new request
	for key, values := range header {
		fmt.Println("Key: ", key)
		for _, value := range values {
			new_req.Header.Add(key, value)
		}
	}

	new_req.Header.Add("Host", req.Host)
	// new_req.Header.Add("Proto", "HTTP/1.1")
	fmt.Println("new request header: ", new_req.Header)
	new_req.Host = req.Host
	fmt.Println("host: ", new_req.Host)
	new_req.Close = true
	new_req.Proto = req.Proto
	new_req.ProtoMajor = 1
	new_req.ProtoMinor = 1

	// Send our modified request to the server and receive the server response
	proxy_client := &http.Client{}
	resp, err := proxy_client.Do(new_req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error sending request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	// Send the server response to the client
	// return the entire response
	resp.Write(w)
}

func handleRequest(conn net.Conn) {
	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("Error receiving client data: %s\n", err)
	}

	separated := strings.Split(request, " ")
	if len(separated) != 3 {
		log.Println("Invalid HTTP Request")
		return
	}

}

// TODO: implement an HTTP proxy
func proxy(server_port string) {
	http.HandleFunc("/", handler)
	fmt.Println("Server listening on " + server_port)

	// ASK: HTTP 1.1 specifies that all HTTP requests must have the Host header explicitly, but assignment assumes there is no such condition?
	http.ListenAndServe(":"+server_port, nil)

	// Create a new HTTP server
	// server := &http.Server{Addr: ":" + server_port}

	// // Listen and serve HTTP requests
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Fatalf("Error with starting proxy server: ", err)
	// }

	// err = server.Shutdown(nil)
	// if err != nil {
	// 	log.Fatalf("Error closing connection: ", err)
	// }
}

// // listen on 127.0.0.1 on port server_port for client connections
// sock, err := net.Listen("tcp", "127.0.0.1:"+server_port)
// if err != nil {
// 	log.Fatalf("Could not setup a socket and listen on port: %s\n", err)
// }
// defer sock.Close()

// for {
// 	// accept any client connection
// 	conn, err := sock.Accept()
// 	if err != nil {
// 		log.Fatalf("Could not connect to client - %s\n", err)
// 	}
// 	// Close connection when this function ends
// 	defer conn.Close()

// 	// setup a reader to continuously read packets from the connection to a buffer so that we can print immediately
// 	// reader w/ buffer of size RECV_BUFFER_SIZE
// 	inp_stream := bufio.NewReaderSize(conn, RECV_BUFFER_SIZE)
// 	buffer := make([]byte, RECV_BUFFER_SIZE)

// 	// for loop handles when entire data is too long to be sent in one go, so it's sent in chunks
// 	for {
// 		// reads bytes from connection (up to RECV_BUFFER_SIZE) and stores them in buffer
// 		num_rec_bytes, err := inp_stream.Read(buffer)

// 		// EOF tells us we read everything
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("Error receiving client data: %s\n", err)
// 		}
// 		// print all the bytes we received as soon as we receive them to stdout
// 		fmt.Print(string(buffer[:num_rec_bytes]))
// 	}
// }

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	proxy(server_port)
}