package cmd

import (
	// "encoding/base64"
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"github.com/urfave/cli"

	"github.com/midoks/dagger/dagger-server/internal/conf"
	"github.com/midoks/dagger/dagger-server/internal/db"
)

var logger *go_logger.Logger

var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts dagger services",
	Description: `Start Http Proxy Server services`,
	Action:      RunService,
	Flags:       []cli.Flag{},
}

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

type SendInfo struct {
	Link        string `json:"link"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	RequestTime string `json:"request_time"`
}

func process(c *gin.Context, ws *websocket.Conn, info *SendInfo) bool {

	// Now, Hijack the writer to get the underlying net.Conn.
	// Which can be either *tcp.Conn, for HTTP, or *tls.Conn, for HTTPS.
	src := ws.UnderlyingConn()
	reader := bufio.NewReader(src)
	defer src.Close()

	dst, err := net.Dial("tcp", info.Link)
	if err != nil {
		logger.Errorf("net.Dial:%s,err:%v", info.Link, err)
		return false
	}
	defer dst.Close()

	// src.SetDeadline(time.Now().Add(5 * time.Second))
	// dst.SetDeadline(time.Now().Add(5 * time.Second))

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

	return true
}

//websocket实现
func websocketReqMethod(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	for {

		mt, message, err := ws.ReadMessage()
		if err != nil {
			logger.Errorf("read websocket msg: %v", err)
			break
		}

		reqInfo := &SendInfo{}
		json.Unmarshal(message, &reqInfo)

		userEnable := conf.GetString("user.enable", "1")
		if userEnable == "1" {

			if db.UserAclCheck(reqInfo.Username, reqInfo.Password) {
				b := process(c, ws, reqInfo)
				if b {
					fmt.Println(reqInfo.Link, "done")
				}
			} else {
				info := fmt.Sprintf("user[%s]:password[%s] acl fail", reqInfo.Username, reqInfo.Password)
				logger.Errorf(info)
				err = ws.WriteMessage(mt, []byte(info))
				if err == nil {
					break
				}
			}

		} else {
			logger.Infof("process[%s]:%d", reqInfo.Link, runtime.NumGoroutine())
			b := process(c, ws, reqInfo)
			if b {
				logger.Infof("process[%s][done]:%d", reqInfo.Link, runtime.NumGoroutine())
			}
		}
	}
}

func RunService(c *cli.Context) error {
	Init()

	httpPort := conf.GetString("http.port", "12345")
	httpPath := conf.GetString("http.path", "ws")

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.GET("/info", func(c *gin.Context) {
		numG := runtime.NumGoroutine()
		c.JSON(200, gin.H{
			"goroutine": numG,
		})
	})

	hp := fmt.Sprintf("/%s", httpPath)
	r.GET(hp, websocketReqMethod)

	r.Run(fmt.Sprintf(":%s", httpPort))
	return nil
}
