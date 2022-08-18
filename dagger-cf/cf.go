package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/midoks/dagger/dagger-cf/internal/cmd"
	"github.com/midoks/dagger/dagger-cf/internal/conf"
)

const Version = "0.0.1"
const AppName = "dagger-cf"

func init() {
	conf.App.Version = Version
	conf.App.Name = AppName
}

func main() {

	app := cli.NewApp()
	app.Name = conf.App.Name
	app.Version = conf.App.Version
	app.Usage = "dagger-cf"
	app.Commands = []cli.Command{
		cmd.Service,
		cmd.Run,
	}

	// fmt.Println(os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Printf("Failed to start application: %v", err)
	}
}
