package task

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/midoks/dagger/dagger-cf/internal/utils"
)

const (
	tcpConnectTimeout = time.Second * 3
	maxRoutine        = 1000
	defaultRoutines   = 200
	defaultPort       = 443
	defaultPingTimes  = 4
	defaultOutput     = "result.csv"
)

var (
	Routines      = defaultRoutines
	TCPPort   int = defaultPort
	PingTimes int = defaultPingTimes
)

const (
	maxDelay = 9999 * time.Millisecond
	minDelay = 0 * time.Millisecond
)

var (
	InputMaxDelay = maxDelay
	InputMinDelay = minDelay
	Output        = defaultOutput
	PrintNum      = 10
)

type Ping struct {
	wg      *sync.WaitGroup
	m       *sync.Mutex
	ips     []*net.IPAddr
	list    PingDelaySet
	control chan bool
	bar     *utils.Bar
}

type PingData struct {
	IP       *net.IPAddr
	Sended   int
	Received int
	Delay    time.Duration
}

type CloudflareIPData struct {
	*PingData
	recvRate      float32
	DownloadSpeed float64
}

// 是否打印测试结果
func NoPrintResult() bool {
	return PrintNum == 0
}

func convertToString(data []CloudflareIPData) [][]string {
	result := make([][]string, 0)
	for _, v := range data {
		result = append(result, v.toString())
	}
	return result
}

// 是否输出到文件
func noOutput() bool {
	return Output == "" || Output == " "
}

func (cf *CloudflareIPData) getRecvRate() float32 {
	if cf.recvRate == 0 {
		pingLost := cf.Sended - cf.Received
		cf.recvRate = float32(pingLost) / float32(cf.Sended)
	}
	return cf.recvRate
}

func (cf *CloudflareIPData) toString() []string {
	result := make([]string, 6)
	result[0] = cf.IP.String()
	result[1] = strconv.Itoa(cf.Sended)
	result[2] = strconv.Itoa(cf.Received)
	result[3] = strconv.FormatFloat(float64(cf.getRecvRate()), 'f', 2, 32)
	result[4] = strconv.FormatFloat(cf.Delay.Seconds()*1000, 'f', 2, 32)
	result[5] = strconv.FormatFloat(cf.DownloadSpeed/1024/1024, 'f', 2, 32)
	return result
}

type PingDelaySet []CloudflareIPData

func (s PingDelaySet) FilterDelay() (data PingDelaySet) {
	if InputMaxDelay > maxDelay || InputMinDelay < minDelay {
		return s
	}
	for _, v := range s {
		if v.Delay > InputMaxDelay { // 平均延迟上限
			break
		}
		if v.Delay < InputMinDelay { // 平均延迟下限
			continue
		}
		data = append(data, v) // 延迟满足条件时，添加到新数组中
	}
	return
}

func (s PingDelaySet) Len() int {
	return len(s)
}

func (s PingDelaySet) Less(i, j int) bool {
	// iRate, jRate := s[i].getRecvRate(), s[j].getRecvRate()
	// if iRate != jRate {
	// 	return iRate < jRate
	// }
	return s[i].Delay < s[j].Delay
}

func (s PingDelaySet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// 下载速度排序
type DownloadSpeedSet []CloudflareIPData

func (s DownloadSpeedSet) Len() int {
	return len(s)
}

func (s DownloadSpeedSet) Less(i, j int) bool {
	return s[i].Received > s[j].Received
}

func (s DownloadSpeedSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DownloadSpeedSet) Print(ipv6 bool) {
	if NoPrintResult() {
		return
	}
	if len(s) <= 0 { // IP数组长度(IP数量) 大于 0 时继续
		fmt.Println("\n[信息] 完整测速结果 IP 数量为 0，跳过输出结果。")
		return
	}
	dateString := convertToString(s) // 转为多维数组 [][]String
	if len(dateString) < PrintNum {  // 如果IP数组长度(IP数量) 小于  打印次数，则次数改为IP数量
		PrintNum = len(dateString)
	}
	headFormat := "%-16s%-5s%-5s%-5s%-6s%-11s\n"
	dataFormat := "%-18s%-8s%-8s%-8s%-10s%-15s\n"
	if ipv6 { // IPv6 太长了，所以需要调整一下间隔
		headFormat = "%-40s%-5s%-5s%-5s%-6s%-11s\n"
		dataFormat = "%-42s%-8s%-8s%-8s%-10s%-15s\n"
	}
	fmt.Printf(headFormat, "IP 地址", "已发送", "已接收", "丢包率", "平均延迟", "下载速度 (MB/s)")
	for i := 0; i < PrintNum; i++ {
		fmt.Printf(dataFormat, dateString[i][0], dateString[i][1], dateString[i][2], dateString[i][3], dateString[i][4], dateString[i][5])
	}
	if !noOutput() {
		fmt.Printf("\n完整测速结果已写入 %v 文件，可使用记事本/表格软件查看。\n", Output)
	}
}

func checkPingDefault() {
	if Routines <= 0 {
		Routines = defaultRoutines
	}
	if TCPPort <= 0 || TCPPort >= 65535 {
		TCPPort = defaultPort
	}
	if PingTimes <= 0 {
		PingTimes = defaultPingTimes
	}
}

func New() *Ping {
	return &Ping{
		wg:      &sync.WaitGroup{},
		m:       &sync.Mutex{},
		control: make(chan bool, Routines),
	}
}

func NewPing() *Ping {
	checkPingDefault()
	ips := loadIPRanges()
	return &Ping{
		wg:      &sync.WaitGroup{},
		m:       &sync.Mutex{},
		ips:     ips,
		list:    make(PingDelaySet, 0),
		control: make(chan bool, Routines),
		bar:     utils.NewBar(len(ips)),
	}
}

func (p *Ping) Run() PingDelaySet {
	if len(p.ips) == 0 {
		return p.list
	}
	// ipVersion := "IPv4"
	// if IPv6 { // IPv6 模式判断
	// 	ipVersion = "IPv6"
	// }
	for _, ip := range p.ips {
		p.wg.Add(1)
		p.control <- false
		go p.start(ip)
	}
	p.wg.Wait()
	sort.Sort(p.list)
	return p.list
}

func (p *Ping) UrlSpeed(url string) time.Duration {
	addr, err := net.ResolveIPAddr("ip", url)
	if err != nil {
		return time.Duration(0)
	}
	recv, totalDlay := p.checkConnection(addr)
	// fmt.Println(recv, totalDlay, url, addr.IP.String())
	if recv == 0 {
		return time.Duration(0)
	}
	delay := totalDlay / time.Duration(recv)
	return time.Duration(int(delay/1000000) * 1000000)
}

func (p *Ping) start(ip *net.IPAddr) {
	defer p.wg.Done()
	p.tcpingHandler(ip)
	<-p.control
}

//bool connectionSucceed float32 time
func (p *Ping) tcping(ip *net.IPAddr) (bool, time.Duration) {
	startTime := time.Now()
	fullAddress := fmt.Sprintf("%s:%d", ip.String(), TCPPort)
	if IPv6 { // IPv6 需要加上 []
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), TCPPort)
	}
	conn, err := net.DialTimeout("tcp", fullAddress, tcpConnectTimeout)
	// fmt.Println(ip.IP.String(), err)
	if err != nil {
		return false, 0
	}
	defer conn.Close()
	duration := time.Since(startTime)
	return true, duration
}

//pingReceived pingTotalTime
func (p *Ping) checkConnection(ip *net.IPAddr) (recv int, totalDelay time.Duration) {
	for i := 0; i < PingTimes; i++ {
		if ok, delay := p.tcping(ip); ok {
			recv++
			totalDelay += delay
		}
	}
	return
}

func (p *Ping) appendIPData(data *PingData) {
	p.m.Lock()
	defer p.m.Unlock()
	p.list = append(p.list, CloudflareIPData{
		PingData: data,
	})
}

// handle tcping
func (p *Ping) tcpingHandler(ip *net.IPAddr) {
	recv, totalDlay := p.checkConnection(ip)
	p.bar.Grow(1)
	if recv == 0 {
		return
	}

	data := &PingData{
		IP:       ip,
		Sended:   PingTimes,
		Received: recv,
		Delay:    totalDlay / time.Duration(recv),
	}
	p.appendIPData(data)
}
