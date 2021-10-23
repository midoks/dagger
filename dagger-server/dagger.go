package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	// "io/ioutil"
	"log"
	"net"
	"net/http"
	// "os"
	"runtime"
	"sync"
	"time"
	// "net/url"
	// "bytes"
	// "strings"
	// "encoding/base64"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func Md5Byte(buf []byte) string {
	hash := md5.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func Md5(s string) string {
	return Md5Byte([]byte(s))
}

//websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SendMsg struct {
	Link    string `json:"link"`
	ReqConn string `json:"reqconn"`
}

var (
	list map[string]net.Conn
)

func init() {
	list = make(map[string]net.Conn)
}

func process(c *gin.Context, ws *websocket.Conn, info *SendMsg) {

	// Now, Hijack the writer to get the underlying net.Conn.
	// Which can be either *tcp.Conn, for HTTP, or *tls.Conn, for HTTPS.
	src := ws.UnderlyingConn()
	reader := bufio.NewReader(src)
	defer src.Close()

	dst, err := net.Dial("tcp", info.Link)
	if err != nil {
		log.Println("net.Dial:", info.Link, err)
		return
	}
	defer dst.Close()

	src.SetDeadline(time.Now().Add(5 * time.Second))
	dst.SetDeadline(time.Now().Add(5 * time.Second))

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
	}()

	wg.Wait()

	fmt.Println("oooo...")
}

//websocket实现
func network(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	for {

		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read ws msg:", err, message)
			break
		}

		res := &SendMsg{}
		json.Unmarshal(message, &res)

		fmt.Println("res:", res.Link)
		fmt.Println("receive", mt, string(message))

		log.Println("process", &c, &ws, runtime.NumGoroutine())
		process(c, ws, res)

		// encodeR := base64.StdEncoding.EncodeToString([]byte(r))
		// d, err := httpGet([]byte("https://www.ixigua.com"))
		// fmt.Println(d, err)
		//写入ws数据
		// err = ws.WriteMessage(mt, []byte(r))
		// if err != nil {
		// 	break
		// }
	}
}

//websocket实现
func ping(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
		err = ws.WriteMessage(mt, []byte("ok"))
		if err != nil {
			break
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/network", network)
	r.GET("/ping", ping)

	r.Run(":12345")
}
