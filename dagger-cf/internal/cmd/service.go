package cmd

// https://github.com/XIU2/CloudflareSpeedTest

import (
	"errors"
	"fmt"
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

func writeHostFile(content string) error {
	file, err := os.OpenFile(defaultHost, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(content))
	return err
}

func writeHost(ip, domain string) error {
	content, err := readHostFile()

	if err != nil {
		return err
	}

	w := fmt.Sprintf("%s\t\t%s\t%s", ip, domain, signKey)
	result := fmt.Sprintf("%s\n%s\n", content, w)
	return writeHostFile(result)
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
	return writeHostFile(w)
}

func RunService(c *cli.Context) error {

	url := c.String("url")
	max_tl := c.Int("max_tl")
	to_host := c.String("to_host")

	task.URL = url
	task.InputMaxDelay = time.Duration(max_tl) * time.Millisecond

	if to_host != "" {
		if url == "" {
			return nil
		}

		if strings.EqualFold(to_host, "yes") {
			t := task.NewPing()
			pingData := t.Run().FilterDelay()
			ip := pingData[0].IP.String()

			err := writeHost(ip, url)
			if err != nil {
				fmt.Println(ip, url, err)
			} else {
				fmt.Println(ip, url)
			}
		}
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
