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
	"time"

	"github.com/gin-gonic/gin"
	go_logger "github.com/phachon/go-logger"
	"github.com/urfave/cli"

	"github.com/midoks/dagger/dagger-server/internal/conf"
)

var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts dagger services",
	Description: `Start Http Proxy Server services`,
	Action:      RunService,
	Flags:       []cli.Flag{},
}

var logger *go_logger.Logger

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

func process(c *gin.Context, ws *websocket.Conn, info *SendInfo) {

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

	fmt.Println("oooo...", info.Link)
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
			logger.Errorf("read websocket msg: %s", err, message)
			break
		}

		reqInfo := &SendInfo{}
		json.Unmarshal(message, &reqInfo)

		fmt.Println("receive", mt, string(message))

		log.Println("process", &c, &ws, runtime.NumGoroutine())

		process(c, ws, reqInfo)

		break

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

func initLogger() {
	logger = go_logger.NewLogger()

	// 文件输出配置
	fileConfig := &go_logger.FileConfig{
		Filename: "./logs/test.log", // 日志输出文件名，不自动存在
		// 如果要将单独的日志分离为文件，请配置LealFrimeNem参数。
		LevelFileName: map[int]string{
			logger.LoggerLevel("error"): "./logs/error.log", // Error 级别日志被写入 error .log 文件
			logger.LoggerLevel("info"):  "./logs/info.log",  // Info 级别日志被写入到 info.log 文件中
			logger.LoggerLevel("debug"): "./logs/debug.log", // Debug 级别日志被写入到 debug.log 文件中
		},
		MaxSize:    1024 * 1024, // 文件最大值（KB），默认值0不限
		MaxLine:    100000,      // 文件最大行数，默认 0 不限制
		DateSlice:  "d",         // 文件根据日期切分， 支持 "Y" (年), "m" (月), "d" (日), "H" (时), 默认 "no"， 不切分
		JsonFormat: false,       // 写入文件的数据是否 json 格式化
		Format:     "",          // 如果写入文件的数据不 json 格式化，自定义日志格式
	}
	// 添加 file 为 logger 的一个输出
	logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)

	logger.Infof("hello,world,now:%s", time.Now().Format("2006/1/2 15:04:05"))
}

func RunService(c *cli.Context) error {
	conf.Load("conf/app.conf")

	httpPort := conf.GetString("http.port", "12345")
	httpPath := conf.GetString("http.path", "ws")

	initLogger()

	r := gin.Default()

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
