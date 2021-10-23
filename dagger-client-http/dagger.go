package main

// https://blog.twofei.com/794/

import (
	"bufio"
	// "bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	listen          = "localhost:8080"
	connectResponse = []byte("HTTP/1.1 200 OK\r\n\r\n")
	username        = "my_username"
	password        = "my_password"
)

var hlock sync.Mutex
var cLink string

func tunnel(w http.ResponseWriter, req *http.Request) {

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
	log.Println("url:", req.RequestURI)
	cLink = req.RequestURI

	// Connect to Remote.
	// dst, err := net.Dial("tcp", req.RequestURI)
	// dst, err := net.Dial("tcp", "dianying.im:80")
	dst, err := net.Dial("tcp", "127.0.0.1:8081")
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

var addr = flag.String("addr", "127.0.0.1:12345", "http service address")

type SendMsg struct {
	Link    string `json:"link"`
	ReqConn string `json:"reqconn"`
}

// TCP Server端测试:处理函数
func process(conn net.Conn) {
	defer conn.Close() // 关闭连接

	for {
		reader := bufio.NewReader(conn)
		src := conn
		// var buf [4096]byte
		// n, err := reader.Read(buf[:]) // 读取数据
		// if err != nil {
		// 	fmt.Println("read from client failed, err: ", err)
		// 	break
		// }
		// recvStr := string(buf[:n])
		fmt.Println("cLink", cLink)
		// fmt.Println("收到Client端发来的数据：", buf[:n])

		// hlock.Unlock()

		//websocket start
		u := url.URL{Scheme: "ws", Host: *addr, Path: "/network"}
		log.Printf("connecting to %s", u.String())

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("dial:", err)
		}
		defer c.Close()

		encodeReq := base64.StdEncoding.EncodeToString([]byte("111"))
		tmp := SendMsg{Link: cLink, ReqConn: encodeReq}

		err = c.WriteJSON(tmp)
		if err != nil {
			log.Println("write:", err)

		}

		// Now, Hijack the writer to get the underlying net.Conn.
		// Which can be either *tcp.Conn, for HTTP, or *tls.Conn, for HTTPS.
		dst := c.UnderlyingConn()
		if err != nil {
			// http.Error(c, err.Error(), http.StatusInternalServerError)
			log.Println("Hijack:", err)
		}
		defer dst.Close()

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			// The returned bufio.Reader may contain unprocessed buffered data from the client.
			// Copy them to dst so we can use src directly.
			if n := reader.Buffered(); n > 0 {
				n64, err := io.CopyN(dst, src, int64(n))
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

			// err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			// if err != nil {
			// 	log.Println("write close:", err)
			// }
		}()

		wg.Wait()

		// go func() {

		// 	for {
		// 		_, message, err := c.ReadMessage()
		// 		if err != nil {
		// 			log.Println("read:", err)
		// 		}
		// 		log.Printf("recv: %s", message)

		// 		// decodeBytes, err := base64.StdEncoding.DecodeString(string(message))
		// 		// if err != nil {
		// 		// 	log.Println("recv base64 decode:", err)

		// 		// }
		// 		// log.Printf("recv ss: %s", string(decodeBytes))

		// 		conn.Write(message) // 发送数据

		// 		break
		// 	}
		// }()

		// err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		// if err != nil {
		// 	log.Println("write close:", err)
		// }

		// websocket end

	}
}

// TCP Server端测试
// 处理函数
func process_local(conn net.Conn) {
	defer conn.Close() // 关闭连接

	for {
		reader := bufio.NewReader(conn)
		// fmt.Println("cLink", cLink)
		// fmt.Println("rev client:", buf[:n])

		// // hlock.Unlock()

		log.Println("local read start")
		dst, err := net.Dial("tcp", cLink)
		if err != nil {
			log.Println("net.Dial:", err)
		}
		defer dst.Close()

		log.Println("local read end")

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()

			// The returned bufio.Reader may contain unprocessed buffered data from the client.
			// Copy them to dst so we can use src directly.
			if n := reader.Buffered(); n > 0 {
				n64, err := io.CopyN(dst, conn, int64(n))
				if n64 != int64(n) || err != nil {
					log.Println("io.CopyN:", n64, err)
					return
				}
			}

			// Relay: src -> dst
			io.Copy(dst, conn)
		}()

		go func() {
			defer wg.Done()

			// Relay: dst -> src
			io.Copy(conn, dst)
		}()

		wg.Wait()
	}
}

func main() {
	go func() {

		ln, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Println(err)
			return
		}
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go process(conn)
		}
	}()

	handler := http.HandlerFunc(tunnel)
	err := http.ListenAndServe(listen, handler)
	if err != http.ErrServerClosed {
		panic(err)
	}
}
