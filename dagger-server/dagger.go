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
	"bytes"
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

	go dealConn(dst, src, info.Link)

	wg.Wait()

	fmt.Println("oooo...")
}

func readAllShut(conn net.Conn) ([]byte, error) { //这个手动方法可以避免粘包的问题
	//bufio.NewWriter
	re := bytes.NewBuffer(nil)
	const N = 666
	for {
		var text [N]byte
		lens, err := conn.Read(text[0:])
		re.Write(text[:lens])
		if lens == 0 || err != nil {
			//log.Println(err)  //在这个死循环里面，不要有任何的输出
			// if errors.As(err,*net.OpError) 	//
			if _, ok := err.(*net.OpError); ok {
				return nil, err
			}
			break
		}
		//conn
		//log.Println(lens,text)
		if lens < N {
			break
		}
	}
	rb := re.Bytes()
	//log.Println(rb,"len",len(rb))
	return rb, nil
	/*data,err:=ioutil.ReadAll(conn)
	if err!=nil{
		log.Printf("读取出现错误%T:%v",err,err)
	}
	return data;*/
}

func dealConn(conn net.Conn, src net.Conn, link string) {

	time.Sleep(10 * time.Second)

	//defer conn.Close()
	//defer conn.Flush()
	//长连接里边的读写操作必须放到循环里面这样才能进行多次的读写
	// 如果连接已经断开，就把这个线程中断掉，怎么判断这个连接已经断开？
	thread_c := 0 //如果连续100秒中读取不到内容，就终止循环
	c := 0
	for {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 666)
				buf = buf[:runtime.Stack(buf, false)]
				log.Printf("运行时错误:%v.Runtime error caught: %s", r, buf)
			}
		}()
		// 注意continue这里也要等待，不然造成内存耗尽，处理器耗尽
		time.Sleep(50 * time.Millisecond)
		//#log.Println(len,string(text))
		thread_c++
		if thread_c > 20*100 {
			log.Println(link, conn.RemoteAddr(), "超过100秒未读取到内容，本连接将关闭")
			conn.Close()
			src.Close()
			c--
			break
		}
		frame, op_err := readAllShut(conn)
		if op_err != nil {
			log.Println(link, conn.RemoteAddr(), "出现读写错误，连接不可用，将会被关闭")
			conn.Close()
			src.Close()
			c--
			break //这种已经关闭的连接，要终止循环，退出这条线程
		}
		if len(frame) == 0 {
			//
			//time.Sleep(50*time.Millisecond)
			continue
		}
		thread_c = 0
		log.Printf("%s:-----------------收到tcp请求:报文的长度是%v,", link, len(frame))
		//TODO
		//这里写自己的业务代码
	}
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
