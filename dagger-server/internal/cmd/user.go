package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/midoks/dagger/dagger-server/internal/db"
)

var User = cli.Command{
	Name:        "user",
	Usage:       "This Command Is User Management",
	Description: `User management [add, delete, modify, query]`,
	Action:      RunUser,
	Flags: []cli.Flag{
		stringFlag("method, m", "query", "add, delete, modify, query"),
		stringFlag("username, u", "", "input username"),
		stringFlag("password, p", "", "input password"),
	},
}

func RunUser(c *cli.Context) error {
	Init()

	argsMethod := c.String("method")
	argsUsername := c.String("username")
	argsPassword := c.String("password")

	if argsMethod == "query" {
		u, _ := db.UsersList()

		if len(u) > 0 {
			fmt.Println("-----------------------------------------------------------------")
			for _, v := range u {
				fmt.Println("|user:", v.Name, "| password:", v.Password, "| status:", v.Status, "|")
			}
			fmt.Println("-----------------------------------------------------------------")
		} else {
			fmt.Println("no has user data!")
		}

	}

	if argsMethod == "add" {
		err := db.UserAdd(argsUsername, argsPassword)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("user add ok!")
		}
	}

	if argsMethod == "delete" {
		err := db.UserDel(argsUsername)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("user delete ok!")
		}
	}

	if argsMethod == "modify" {
		err := db.UserMod(argsUsername, argsPassword)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("user modify ok!")
		}
	}

	return nil
}
