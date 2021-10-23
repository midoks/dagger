package main

import (
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	// "net/url"
	// "bytes"
	// "strings"
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func httpGet(url []byte) ([]byte, error) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", string(url), nil) //建立一个请求
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
	//Add 头协议
	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Cookie", "设置cookie")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	response, err := client.Do(reqest) //提交
	defer response.Body.Close()
	// cookies := response.Cookies() //遍历cookies
	// for _, cookie := range cookies {
	// fmt.Println("cookie:", cookie)
	// }

	body, err := ioutil.ReadAll(response.Body)

	return body, err
}

type SendMsg struct {
	Link    string `json:"link"`
	ReqConn string `json:"reqconn"`
}

func process(c *gin.Context, ws *websocket.Conn, info *SendMsg) string {

	// decodeBytes, err := base64.StdEncoding.DecodeString(info.ReqConn)
	// if err != nil {
	// 	log.Println("decodeBytes:", err)
	// }

	// fmt.Println(decodeBytes)
	dst, err := net.Dial("tcp", info.Link)
	if err != nil {
		log.Println("net.Dial:", err)
	}
	defer dst.Close()

	// Now, Hijack the writer to get the underlying net.Conn.
	// Which can be either *tcp.Conn, for HTTP, or *tls.Conn, for HTTPS.
	src := ws.UnderlyingConn()
	reader := bufio.NewReader(src)
	if err != nil {
		// http.Error(c, err.Error(), http.StatusInternalServerError)
		log.Println("Hijack:", err)
		return "ee"
	}
	defer src.Close()

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

	return "buf.String()"
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
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read", err)
			break
		}

		res := &SendMsg{}
		json.Unmarshal(message, &res)
		fmt.Println("res:", res.Link)

		fmt.Println("receive", mt, string(message))

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
