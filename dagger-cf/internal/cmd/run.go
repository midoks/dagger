package cmd

import (
	"fmt"

	"github.com/midoks/dagger/dagger-cf/internal/task"
	"github.com/urfave/cli"
)

var Run = cli.Command{
	Name:        "run",
	Usage:       "This command run services",
	Description: `Start Station speed measurement Services`,
	Action:      RunRun,
	Flags: []cli.Flag{
		stringFlag("url, u", "", "Custom Configuration Url"),
	},
}

func RunRun(c *cli.Context) error {
	url := c.String("url")
	delay := task.New().UrlSpeed(url)
	fmt.Println(delay)
	return nil
}
