/*****************************************************************************
 * http_proxy_DNS.go
 * Names: Sam Liang
 * NetIds: saml
 *****************************************************************************/

// TODO: implement an HTTP proxy with DNS Prefetching

// Note: it is highly recommended to complete http_proxy.go first, then copy it
// with the name http_proxy_DNS.go, thus overwriting this file, then edit it
// to add DNS prefetching (don't forget to change the filename in the header
// to http_proxy_DNS.go in the copy of http_proxy.go)

package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func fetchDNS(doc *html.Node) {
	nodes := []*html.Node{}
	nodes = append(nodes, doc)

	for len(nodes) > 0 {
		// get the newest node (dfs) and pop it
		num_nodes := len(nodes)
		node := nodes[num_nodes-1]
		nodes = nodes[:num_nodes-1]

		//
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					net.LookupHost(attr.Val)
					break
				}
			}
		}

		// populate with next nodes
		next := node.FirstChild
		for next != nil {
			nodes = append(nodes, next)
			next = next.NextSibling
		}
	}
}

func return_error(conn net.Conn) {
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

func handleRequest(conn net.Conn, reader *bufio.Reader) {
	req, err := http.ReadRequest(reader)
	if err != nil {
		return_error(conn)
		return
	}

	if req.Method != "GET" {
		// create response with 500 error
		return_error(conn)
		return
	}

	// new_req.URL, _ = url.Parse(req.URL.Path)

	// not allowed to have RequestURI set when doing a client request.
	// req.RequestURI = ""

	// relative URL
	// req.URL.Path = req.URL.String()[strings.Index(req.URL.String()[8:], "/")+8:]
	// req.Close = true // new_req.Header.Set("Connection", "close")
	req.Header.Set("Connection", "close")
	// // new_req.Header.Add("Scheme", req.URL.Scheme)
	// // new_req.Header.Add("Proto", "HTTP/1.1")

	// new_req.Close = true // new_req.Header.Set("Connection", "close")
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1

	tsp := &http.Transport{}
	resp, err := tsp.RoundTrip(req)

	if err != nil {
		return_error(conn)
		return
	}

	// DNS fetching
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return_error(conn)
		return
	}

	graph, err := html.Parse(bytes.NewReader(body))

	if err != nil {
		return_error(conn)
		return
	}

	go fetchDNS(graph)

	// return the entire response to the client
	_, err = conn.Write(body)

	if err != nil {
		return_error(conn)
		return
	}

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
			log.Printf("Could not connect to client - %s\n", err)
			continue
		}

		// setup a reader to continuously read packets from the connection to a buffer so that we can print immediately
		// reader w/ buffer of size RECV_BUFFER_SIZE
		inp_stream := bufio.NewReader(conn)
		go handleRequest(conn, inp_stream)
	}
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]
	proxy(server_port)
}
