package cmd

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
	"github.com/valyala/fastjson"

	"github.com/midoks/dagger/dagger-server/internal/conf"
	"github.com/midoks/dagger/dagger-server/internal/db"
	"github.com/midoks/dagger/dagger-server/internal/debug"
	"github.com/midoks/dagger/dagger-server/internal/log"
)

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

// Byte to string, only read-only
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//String to byte, only read-only
func StringToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

//websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Simple Process
func process(c *gin.Context, mt int, ws *websocket.Conn, link string) bool {

	dst, err := net.DialTimeout("tcp", link, time.Second)
	if err != nil {
		info := fmt.Sprintf("net.DialTimeout: %s, error: %v", link, err)
		log.Errorf(info)
		ws.WriteMessage(mt, StringToBytes(info))
		return false
	}

	defer dst.Close()

	src := ws.UnderlyingConn()
	reader := bufio.NewReader(src)
	defer src.Close()

	src.SetDeadline(time.Now().Add(conf.Http.Timeout * time.Second))
	dst.SetDeadline(time.Now().Add(conf.Http.Timeout * time.Second))

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		// The returned bufio.Reader may contain unprocessed buffered data from the client.
		// Copy them to dst so we can use src directly.
		if n := reader.Buffered(); n > 0 {
			n64, err := io.CopyN(dst, src, int64(n))
			if n64 != int64(n) || err != nil {
				log.Errorf("io.CopyN:%d, err:%T", n64, err)
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

//Simple Process
func process2(c *gin.Context, src net.Conn, reader *bufio.Reader, link string) bool {

	dst, err := net.Dial("tcp", link)
	if err != nil {
		log.Errorf("net.Dial:%s,err:%v", link, err)
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
				log.Errorf("io.CopyN:%d, err:%T", n64, err)
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

func checkConn(conn net.Conn) (net.Conn, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	var one = []byte{}
	_, err := conn.Read(one)
	if err != nil {
		return conn, err
	}
	if err == io.EOF {
		return conn, err
	}
	var zero time.Time
	conn.SetReadDeadline(zero)
	return conn, nil
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
			log.Errorf("read websocket msg error: %v:%d", err, runtime.NumGoroutine())
			break
		}

		var p fastjson.Parser

		reqMessage := BytesToString(message)
		v, err := p.Parse(reqMessage)
		if err != nil {
			log.Errorf("cannot parse json: %s", err)
			break
		}

		link := BytesToString(v.GetStringBytes("link"))
		username := BytesToString(v.GetStringBytes("username"))
		password := BytesToString(v.GetStringBytes("password"))

		log.Infof("P[%s]:%d", link, runtime.NumGoroutine())
		startTime := time.Now()
		if conf.User.Enable {
			if !db.UserAclCheck(username, password) {
				info := fmt.Sprintf("user[%s]:password[%s] acl fail", username, password)
				log.Errorf(info)
				ws.WriteMessage(mt, []byte(info))
				break
			}
		}

		b := process(c, mt, ws, link)
		tcTime := time.Since(startTime)
		if b {
			log.Infof("P[%s][done][%v]:%d", link, tcTime, runtime.NumGoroutine())
		} else {
			log.Errorf("P[%s][fali]:%d", link, runtime.NumGoroutine())
		}
		break
	}

}

func RunService(c *cli.Context) error {

	httpPort := conf.Http.Port
	httpPath := conf.Http.Path

	r := gin.Default()
	// r.SetMode(gin.ReleaseMode)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.GET("/info", func(c *gin.Context) {
		numG := runtime.NumGoroutine()
		c.JSON(200, gin.H{
			"goroutine": numG,
		})
	})

	runMode := conf.App.RunMode
	if strings.EqualFold(runMode, "dev") {
		go debug.Pprof()
	}

	hp := fmt.Sprintf("/%s", httpPath)
	r.GET(hp, websocketReqMethod)

	r.Run(fmt.Sprintf(":%d", httpPort))
	return nil
}
