package cmd

import (
	"fmt"
	// "os"
	// "path/filepath"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
)

var defaultClientsCount = runtime.NumCPU()

var (
	fakeResponse = []byte("Hello, world!")
	getRequest   = "GET /foobar?baz HTTP/1.1\r\nHost: google.com\r\nUser-Agent: aaa/bbb/ccc/ddd/eee Firefox Chrome MSIE Opera\r\n" +
		"Referer: http://example.com/aaa?bbb=ccc\r\nCookie: foo=bar; baz=baraz; aa=aakslsdweriwereowriewroire\r\n\r\n"
	postRequest = fmt.Sprintf("POST /foobar?baz HTTP/1.1\r\nHost: google.com\r\nContent-Type: foo/bar\r\nContent-Length: %d\r\n"+
		"User-Agent: Opera Chrome MSIE Firefox and other/1.2.34\r\nReferer: http://google.com/aaaa/bbb/ccc\r\n"+
		"Cookie: foo=bar; baz=baraz; aa=aakslsdweriwereowriewroire\r\n\r\n%q",
		len(fakeResponse), fakeResponse)
)

type RequestContext func(ctx *gin.Context)

func init() {
	fmt.Println("mock")
}

// go test -bench=BenchmarkServerWS -benchmem -benchtime=1s -memprofile mem.out
// go test -bench=BenchmarkServerWS -benchmem -benchtime=1s
func BenchmarkServerWS(b *testing.B) {
	router := gin.New()
	router.Any("/network", websocketReqMethod)
	runRequest(b, router, "GET", "/network")
}

// go test -bench=BenchmarkServerDebug -benchmem -benchtime=1s
func BenchmarkServerDebug(b *testing.B) {
	router := gin.New()
	router.Any("/ping", func(ctx *gin.Context) {})
	runRequest(b, router, "GET", "/ping")
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func runRequest(B *testing.B, r *gin.Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	// w := newMockWriter()
	w := httptest.NewRecorder()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}
