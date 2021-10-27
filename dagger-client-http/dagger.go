package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/midoks/dagger/dagger-client-http/internal/cmd"
	"github.com/midoks/dagger/dagger-client-http/internal/conf"
)

const Version = "0.0.1"
const AppName = "dagger-client-http"

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
	}

	fmt.Println(os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Printf("Failed to start application: %v", err)
	}
}
