package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/midoks/dagger/dagger-server/internal/db"
)

var User = cli.Command{
	Name:        "user",
	Usage:       "This Command Is User Management",
	Description: `User management [add, delete, modify, query, list]`,
	Action:      RunUser,
	Flags: []cli.Flag{
		stringFlag("method, m", "list", "add, delete, modify, query,list"),
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

		if argsUsername == "" {
			fmt.Println("user cannot be empty!")
			return nil
		}

		u, err := db.UserGetByName(argsUsername)

		if err == nil {
			fmt.Println("-----------------------------------------------------------------")
			fmt.Println("|user:", u.Name, "| password:", u.Password, "| status:", u.Status, "|")
			fmt.Println("-----------------------------------------------------------------")
		} else {
			fmt.Println("no has user[", argsUsername, "] data!")
		}

	}

	if argsMethod == "list" {
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

		if argsUsername == "" {
			fmt.Println("user cannot be empty!")
			return nil
		}

		if argsPassword == "" {
			fmt.Println("user cannot be empty!")
			return nil
		}

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

		if argsUsername == "" {
			fmt.Println("user cannot be empty!")
			return nil
		}

		if argsPassword == "" {
			fmt.Println("user cannot be empty!")
			return nil
		}

		err := db.UserMod(argsUsername, argsPassword)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("user modify ok!")
		}
	}

	return nil
}
