package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pepper"
	app.Version = "0.1.0"
	app.Usage = "pepper <target> <function> [ARGUMENTS ...]"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "hostname, H",
			Usage:  "Salt API hostname. Should include http[s]//.",
			EnvVar: "SALT_HOST",
		},
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "Salt API username.",
			EnvVar: "SALT_USER",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "Salt API password.",
			EnvVar: "SALT_PASSWORD",
		},

		cli.StringFlag{
			Name:   "auth, a",
			Value:  "ldap",
			Usage:  "Salt authentication method.",
			EnvVar: "SALT_AUTH",
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) < 2 {
			fmt.Println("pepper <target> <function> [ARGUMENTS ...]")
			return
		}

		hostname := c.String("hostname")
		username := c.String("username")
		password := c.String("password")
		auth := c.String("auth")

		salt := NewSalt(hostname)

		err := salt.Login(username, password, auth)
		if err != nil {
			log.Fatal(err)
		}

		target := c.Args().Get(0)
		function := c.Args().Get(1)
		arguments := c.Args().Get(2)

		response, _ := salt.Run(target, function, arguments)
		fmt.Println(response)
	}

	app.Run(os.Args)

}
