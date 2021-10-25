package main

// https://blog.twofei.com/794/

import (
	// "bufio"
	// "bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	// "net/url"
	"runtime"
	"sync"
	// "time"

	"github.com/gorilla/websocket"
)

var (
	listen          = "localhost:1097"
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
	log.Println("url:", req.RequestURI, runtime.NumGoroutine())

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

type RunServer struct {
}

var wsConn *websocket.Conn
var err error

func tunnel_ws(w http.ResponseWriter, req *http.Request) {

	// We handle CONNECT method only
	if req.Method != http.MethodConnect {
		log.Println(req.Method, req.RequestURI)
		http.NotFound(w, req)
		return
	}
	// hlock.Lock()

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
	log.Println("url:", req.RequestURI, runtime.NumGoroutine())

	//websocket start
	// u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	// log.Printf("connecting to %s", u.String())

	link := "ws://127.0.0.1:12345/network"
	// link := "wss://v3.biqu.xyz/ws"

	// if wsConn == nil {
	wsConn, _, err = websocket.DefaultDialer.Dial(link, nil)
	if err != nil {
		log.Println("ws dial:", err)
		return
	}

	encodeReq := base64.StdEncoding.EncodeToString([]byte("111"))
	tmp := SendMsg{Link: req.RequestURI, ReqConn: encodeReq}

	err := wsConn.WriteJSON(tmp)
	if err != nil {
		log.Println("write:", err)

		wsConn, _, err = websocket.DefaultDialer.Dial(link, nil)
		if err != nil {
			log.Println("ws dial2:", err)
			return
		}
		fmt.Println("wsConn3:", &wsConn)

		err := wsConn.WriteJSON(tmp)
		if err != nil {
			return
		}
	}

	// Connect to Remote.
	dst := wsConn.UnderlyingConn()
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

	// src.SetDeadline(time.Now().Add(5 * time.Second))
	// dst.SetDeadline(time.Now().Add(5 * time.Second))

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

// var addr = flag.String("addr", "127.0.0.1:12345", "http service address")
var addr = flag.String("addr", "v2.biqu.xyz", "http service address")

type SendMsg struct {
	Link    string `json:"link"`
	ReqConn string `json:"reqconn"`
}

func main() {

	handler := http.HandlerFunc(tunnel_ws)
	err := http.ListenAndServe(listen, handler)
	if err != http.ErrServerClosed {
		panic(err)
	}
}
