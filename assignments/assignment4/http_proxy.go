/*****************************************************************************
 * http_proxy.go
 * Names: Sam Liang
 * NetIds: saml
 *****************************************************************************/

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// fmt.Println("Absolute URL: ", req.URL.String())
	// relative_url := req.URL.Path
	// fmt.Println("relative url: ", req.URL.Path)

	// this URL should be a relative URL and not the absolute URL
	new_req, err := http.NewRequest("GET", req.URL.Path, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error creating request: ", err)
		return
	}

	// new_req.URL, _ = url.Parse(req.URL.Path)

	// Set the header values in the new request
	for key, values := range req.Header {
		for _, value := range values {
			new_req.Header.Add(key, value)
		}
	}

	// new_req.Header.Add("Host", req.Host)
	// new_req.Header.Add("Scheme", req.URL.Scheme)
	// new_req.Header.Add("Proto", "HTTP/1.1")
	// fmt.Println("new request header: ", new_req.Header)

	new_req.URL.Scheme = req.URL.Scheme
	new_req.URL.Host = req.Host
	// new_req.Host = req.Host
	new_req.Close = true
	new_req.Proto = req.Proto

	// reqDump, err := httputil.DumpRequest(new_req, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("REQUEST:\n%s", string(reqDump))
	// fmt.Println()

	// Send our modified request to the server and receive the server response
	proxy_client := &http.Client{}
	resp, err := proxy_client.Do(new_req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error sending request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := resp.Body.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
		}
	}

	// resp.Body.Read(body)

	// body, _ := io.ReadAll(resp.Body)

	// Set the Connection header to "close".
	w.Header().Set("Connection", "close")

	// Send the server response to the client
	// return the entire response\
	w.Write(b)
	// resp.Write(body)
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

	// ASK: HTTP 1.1 specifies that all HTTP requests must have the Host header explicitly, but assignment assumes there is no such condition?
	// http.ListenAndServe(":"+server_port, nil)

	// Create a new HTTP server
	server := &http.Server{Addr: ":" + server_port}

	// Listen and serve HTTP requests
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error with starting proxy server: ", err)
	}
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
