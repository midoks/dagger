package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/midoks/dagger/dagger-server/internal/cmd"
	"github.com/midoks/dagger/dagger-server/internal/conf"
)

const Version = "0.0.5"
const AppName = "dagger-server"

func init() {
	conf.App.Version = Version
	conf.App.Name = AppName
}

func main() {

	app := cli.NewApp()
	app.Name = conf.App.Name
	app.Version = conf.App.Version
	app.Usage = "A simple http proxy service"
	app.Commands = []cli.Command{
		cmd.Service,
		cmd.User,
	}

	if err := app.Run(os.Args); err != nil {
		log.Printf("Failed to start application: %v", err)
	}
}
