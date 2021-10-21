package main

// https://blog.twofei.com/794/

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

var (
	listen          = "localhost:8080"
	connectResponse = []byte("HTTP/1.1 200 OK\r\n\r\n")
	username        = "my_username"
	password        = "my_password"
)

func tunnel(w http.ResponseWriter, req *http.Request) {
	// We handle CONNECT method only
	if req.Method != http.MethodConnect {
		log.Println(req.Method, req.RequestURI)
		http.NotFound(w, req)
		return
	}

	// Proxy-Authorization is set by client software.
	// Authorization is used by req.BasicAuth().

	// req.Header.Set("Authorization", req.Header.Get("Proxy-Authorization"))
	// user, pass, ok := req.BasicAuth()
	// if !ok || !(user == username && pass == password) {
	// 	log.Println("bad credential.", "user:", user, "pass:", pass)
	// 	// Don't let them know we support CONNECT.
	// 	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// The host:port pair.
	log.Println("req:", req.Method, req)
	log.Println("url:", req.RequestURI)

	// Connect to Remote.
	dst, err := net.Dial("tcp", req.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer dst.Close()

	// Upon success, we respond a 200 status code to client.
	w.Write(connectResponse)

	// Now, Hijack the writer to get the underlying net.Conn.
	// Which can be either *tcp.Conn, for HTTP, or *tls.Conn, for HTTPS.
	src, bio, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer src.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		// The returned bufio.Reader may contain unprocessed buffered data from the client.
		// Copy them to dst so we can use src directly.
		if n := bio.Reader.Buffered(); n > 0 {
			n64, err := io.CopyN(dst, bio, int64(n))
			if n64 != int64(n) || err != nil {
				log.Println("io.CopyN:", n64, err)
				return
			}
		}

		// Relay: src -> dst
		io.Copy(dst, src)
	}()

	go func() {
		defer wg.Done()

		// Relay: dst -> src
		io.Copy(src, dst)
	}()

	wg.Wait()
}

func main() {
	handler := http.HandlerFunc(tunnel)
	err := http.ListenAndServe(listen, handler)
	if err != http.ErrServerClosed {
		panic(err)
	}
}
