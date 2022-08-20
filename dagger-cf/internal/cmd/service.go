package cmd

// https://github.com/XIU2/CloudflareSpeedTest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/midoks/dagger/dagger-cf/internal/task"
	"github.com/urfave/cli"
)

var Service = cli.Command{
	Name:        "service",
	Usage:       "This command starts services",
	Description: `Start Cloudflare IP Preference Services`,
	Action:      RunService,
	Flags: []cli.Flag{
		stringFlag("url, u", "", "Custom Configuration URL"),
		stringFlag("ipv4, v4", "", "Custom Configuration IPV4"),
		stringFlag("to_host, th", "no", "Custom Configuration Set Host yes|no|clean, default:no"),
		intFlag("max_tl, ml", 200, "Average delay upper limit; Only output IP with lower than the specified average delay, which can be matched with other upper / lower limits; (default 9999 MS)"),
	},
}

var (
	defaultHost = "/etc/hosts"
	signKey     = "#Dagger Hosts Don`t Remove and Change"
)

func replaceStringByRegex(str, rule, replace string) (string, error) {
	reg, err := regexp.Compile(rule)
	if reg == nil || err != nil {
		return "", errors.New("replaceStringByRegex error:" + err.Error())
	}
	return reg.ReplaceAllString(str, replace), nil
}

func readHostFile() (string, error) {
	content, err := os.ReadFile(defaultHost)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func writeHostAppendFile(content string) error {
	file, err := os.OpenFile(defaultHost, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(content))
	return err
}

func writeHostFile(content string) error {
	ioutil.WriteFile(defaultHost, []byte(content), os.ModePerm)
	return nil
}

func writeHost(ip, domain string) error {
	w := fmt.Sprintf("\n%s\t\t%s\t%s", ip, domain, signKey)
	return writeHostFile(w)
}

func writeHostAppend(ip, domain string) error {
	w := fmt.Sprintf("\n%s\t\t%s\t%s", ip, domain, signKey)
	return writeHostAppendFile(w)
}

func writeHostAppendN() error {
	return writeHostAppendFile("\n")
}

func clearHost() error {
	content, err := readHostFile()
	if err != nil {
		return err
	}

	w, err := replaceStringByRegex(content, ".*"+signKey, "")
	if err != nil {
		return err
	}

	w = strings.TrimSpace(w)

	fmt.Println(w)
	return writeHostFile(w)
}

func RunService(c *cli.Context) error {

	url := c.String("url")
	max_tl := c.Int("max_tl")
	to_host := c.String("to_host")
	ipv4 := c.String("ipv4")

	task.URL = url
	task.InputMaxDelay = time.Duration(max_tl) * time.Millisecond
	if ipv4 != "" {
		task.IPFile = ipv4
	}

	if url != "" && strings.EqualFold(to_host, "yes") {
		t := task.NewPing()
		pingData := t.Run().FilterDelay()

		speedData := task.TestDownloadSpeed(pingData)
		speedData.Print(task.IPv6)

		ip := speedData[0].IP.String()
		// ip := "127.0.0.1"
		url = strings.Trim(url, "\"")
		url = strings.Trim(url, "|")
		hlist := strings.Split(url, "|")
		clearHost()
		for _, url_val := range hlist {
			err := writeHostAppend(ip, url_val)
			if err != nil {
				fmt.Println(ip, url_val, err)
			} else {
				fmt.Println(ip, url_val)
			}
		}
		writeHostAppendN()
	}

	if strings.EqualFold(to_host, "clean") {
		err := clearHost()
		if err != nil {
			fmt.Println("clear host:", err)
		} else {
			fmt.Println("clear host: ok")
		}

	}

	// for _, i := range pingData {
	// 	fmt.Println(i.IP.String(), i.Delay)
	// }

	return nil
}
