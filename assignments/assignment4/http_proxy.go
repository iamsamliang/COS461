/*****************************************************************************
 * http_proxy.go
 * Names: Sam Liang
 * NetIds: saml
 *****************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func handler(w http.ResponseWriter, req *http.Request) {
	// if req.Method != "GET" {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// // fmt.Println("Absolute URL: ", req.URL.String())
	// // relative_url := req.URL.Path
	// // fmt.Println("relative url: ", req.URL.Path)

	// // this URL should be a relative URL and not the absolute URL
	// new_req, err := http.NewRequest("GET", req.URL.Path, nil)
	// if err != nil {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	log.Fatalf("Error creating request: ", err)
	// 	// need to create response with 500 error then do the same thing.
	// 	resp.Write(body) // resp.Write(to socket)
	// 	resp.Body.Close()
	// 	net.Close
	// 	return
	// }

	// // new_req.URL, _ = url.Parse(req.URL.Path)

	// // Set the header values in the new request
	// for key, values := range req.Header {
	// 	for _, value := range values {
	// 		new_req.Header.Add(key, value)
	// 	}
	// }

	// // new_req.Header.Add("Host", req.Host)
	// // new_req.Header.Add("Scheme", req.URL.Scheme)
	// // new_req.Header.Add("Proto", "HTTP/1.1")
	// // fmt.Println("new request header: ", new_req.Header)

	// new_req.URL.Scheme = req.URL.Scheme
	// new_req.URL.Host = req.Host
	// // new_req.Host = req.Host
	// new_req.Close = true // new_req.Header.Set("Connection", "close")
	// new_req.Proto = req.Proto

	// // reqDump, err := httputil.DumpRequest(new_req, true)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// // fmt.Printf("REQUEST:\n%s", string(reqDump))
	// // fmt.Println()

	// // Send our modified request to the server and receive the server response
	// proxy_client := &http.Client{}
	// resp, err := proxy_client.Do(new_req)
	// if err != nil {
	// 	// return a response with error
	// 	// if err is EOF I still return 500 error
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	log.Fatalf("Error sending request: %s\n", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// b := make([]byte, 0, 512)
	// for {
	// 	if len(b) == cap(b) {
	// 		// Add more capacity (let append pick how much).
	// 		b = append(b, 0)[:len(b)]
	// 	}
	// 	n, err := resp.Body.Read(b[len(b):cap(b)])
	// 	b = b[:len(b)+n]
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			err = nil
	// 			break
	// 		}
	// 	}
	// }

	// resp.Body.Read(body)

	// body, _ := io.ReadAll(resp.Body)

	// Set the Connection header to "close".
	// w.Header().Set("Connection", "close")

	// Send the server response to the client
	// return the entire response\
	// w.Write(b)
	// resp.Write(body) // resp.Write(to socket)
	// resp.Body.Close()
	// net.Close
}

func handleRequest(conn net.Conn, reader *bufio.Reader) {
	req, err := http.ReadRequest(reader)
	if err != nil {
		// return a response with error
		// if err is EOF I still return 500 error
		// create response with 500 error
		resp := &http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: 500,
			Body:       ioutil.NopCloser(strings.NewReader("")),
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		resp.Write(conn)
		resp.Body.Close()
		conn.Close()
		return
	}

	if req.Method != "GET" {
		// create response with 500 error
		resp := &http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: 500,
			Body:       ioutil.NopCloser(strings.NewReader("")),
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		resp.Write(conn)
		resp.Body.Close()
		conn.Close()
		return
	}

	// new_req.URL, _ = url.Parse(req.URL.Path)

	// not allowed to have RequestURI set when doing a client request.
	// req.RequestURI = ""

	// relative URL
	// req.URL.Path = req.URL.String()[strings.Index(req.URL.String()[8:], "/")+8:]
	// req.URL, _ = url.Parse(req.URL.Path)
	// req.Close = true // new_req.Header.Set("Connection", "close")
	req.Header.Set("Connection", "close")
	// // new_req.Header.Add("Scheme", req.URL.Scheme)
	// // new_req.Header.Add("Proto", "HTTP/1.1")

	// new_req.Close = true // new_req.Header.Set("Connection", "close")
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1

	// Send our modified request to the server and receive the server response
	// proxy_client := &http.Client{}
	// resp, err := proxy_client.Do(req)

	// Connect to the server w/ server_ip at server_port
	// server_conn, err := net.Dial("tcp", req.URL.String()+":80")
	// if err != nil {
	// 	// create response with 500 error
	// 	resp := &http.Response{
	// 		Status:     "500 Internal Server Error",
	// 		StatusCode: 500,
	// 		Body:       ioutil.NopCloser(strings.NewReader("")),
	// 		ProtoMajor: 1,
	// 		ProtoMinor: 1,
	// 	}
	// 	resp.Write(conn)
	// 	resp.Body.Close()
	// 	conn.Close()
	// 	return
	// }

	// maybe create an io.Writer via server_conn
	// req.Write(server_conn)

	// server_conn.Read()

	// create a buffer to store the data from stdin so we can immediately send it to server
	// create a reader to read from stdin and put the data in buffer. Read at most SEND_BUFFER_SIZE
	// buffer := make([]byte, SEND_BUFFER_SIZE)
	// inp_stream := bufio.NewReaderSize(os.Stdin, SEND_BUFFER_SIZE)

	tsp := &http.Transport{}
	resp, err := tsp.RoundTrip(req)

	if err != nil {
		// return a response with error
		// if err is EOF I still return 500 error
		// create response with 500 error
		fmt.Println("proxy_client.Do error: ", err)
		resp := &http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: 500,
			Body:       ioutil.NopCloser(strings.NewReader("")),
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		resp.Write(conn)
		resp.Body.Close()
		conn.Close()
		return
	}

	// resp.Body.Read(body)

	// body, _ := io.ReadAll(resp.Body)

	// Set the Connection header to "close".
	// w.Header().Set("Connection", "close")

	// return the entire response to the client
	resp.Write(conn) // resp.Write(to socket)
	resp.Body.Close()
	conn.Close()
}

// TODO: implement an HTTP proxy
func proxy(server_port string) {
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

		// setup a reader to continuously read packets from the connection to a buffer so that we can print immediately
		// reader w/ buffer of size RECV_BUFFER_SIZE
		inp_stream := bufio.NewReader(conn)
		go handleRequest(conn, inp_stream)

		// buffer := make([]byte, 4096)

		// for loop handles when entire data is too long to be sent in one go, so it's sent in chunks
		// for {
		// 	// reads bytes from connection (up to RECV_BUFFER_SIZE) and stores them in buffer
		// 	num_rec_bytes, err := inp_stream.Read(buffer)

		// 	// EOF tells us we read everything
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	if err != nil {
		// 		log.Fatalf("Error receiving client data: %s\n", err)
		// 	}
		// 	// print all the bytes we received as soon as we receive them to stdout
		// 	fmt.Print(string(buffer[:num_rec_bytes]))
		// }

		// Close connection when this function ends
		// conn.Close()
	}

	// just using HTTP code below
	// http.HandleFunc("/", handler)

	// ASK: HTTP 1.1 specifies that all HTTP requests must have the Host header explicitly, but assignment assumes there is no such condition?
	// http.ListenAndServe(":"+server_port, nil)

	// Create a new HTTP server
	// server := &http.Server{Addr: ":" + server_port}

	// Listen and serve HTTP requests
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Fatalf("Error with starting proxy server: ", err)
	// }
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	proxy(server_port)
}
