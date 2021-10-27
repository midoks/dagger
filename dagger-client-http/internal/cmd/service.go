package cmd

// https://blog.twofei.com/794/

import (
	// "encoding/base64"
	// "fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts services",
	Description: `Start Http Proxy services`,
	Action:      RunService,
	Flags: []cli.Flag{
		stringFlag("port, p", "localhost:1097", "Custom Configuration Port"),
		stringFlag("websocket, w", "", "Custom Configuration WebSocket"),
		stringFlag("username, u", "", "Custom Configuration Username"),
		stringFlag("password, m", "", "Custom Configuration Password"),
	},
}

var (
	listen          = "localhost:1097"
	connectResponse = []byte("HTTP/1.1 200 OK\r\n\r\n")
	username        string
	password        string
	websocketLink   string
	wsConn          *websocket.Conn
	err             error
)

type SendInfo struct {
	Link        string `json:"link"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	RequestTime string `json:"request_time"`
}

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

func tunnelWs(w http.ResponseWriter, req *http.Request) {

	// We handle CONNECT method only
	if req.Method != http.MethodConnect {
		log.Println(req.Method, req.RequestURI)
		http.NotFound(w, req)
		return
	}

	// The host:port pair.
	log.Println("url:", req.RequestURI, runtime.NumGoroutine())

	// link := "ws://127.0.0.1:12345/network"
	// link := "wss://v3.biqu.xyz/ws"

	// fmt.Println(link)
	// fmt.Println("websocketLink:", websocketLink)
	// fmt.Println("username:", username)
	// fmt.Println("password:", password)
	// fmt.Println("listen:", listen)

	wsConn, _, err = websocket.DefaultDialer.Dial(websocketLink, nil)
	if err != nil {
		log.Println("ws dial:", err)
		return
	}

	tmp := SendInfo{
		Link:        req.RequestURI,
		Username:    username,
		Password:    password,
		RequestTime: time.Now().Format("2006/1/2 15:04:05"),
	}
	// fmt.Println(tmp)

	err := wsConn.WriteJSON(tmp)
	if err != nil {
		log.Println("write:", err)
		return
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

func RunService(c *cli.Context) error {

	username = c.String("username")
	password = c.String("password")
	listen = c.String("port")
	websocketLink = c.String("websocket")

	if websocketLink == "" {

		//本地监听 -> 本地请求数据
		handler := http.HandlerFunc(tunnel)
		err := http.ListenAndServe(listen, handler)
		if err != http.ErrServerClosed {
			return err
		}
	} else {

		//本地监听 -> websocket请求数据
		handler := http.HandlerFunc(tunnelWs)
		err := http.ListenAndServe(listen, handler)
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
