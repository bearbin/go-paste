package main

import (
	"github.com/bearbin/go-paste/pastebin"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "get and put pastes from pastebin and other paste sites."
	app.Commands = []cli.Command{
		{
			Name:      "put",
			ShortName: "p",
			Usage:     "get a username from a uuid",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "title, t", Value: "", Usage: "the title for the paste"},
			},
			Action: func(c *cli.Context) {
				var err error
				var text []byte
				if c.Args().First() == "-" {
					text, err = ioutil.ReadAll(os.Stdin)
				} else {
					text, err = ioutil.ReadFile(c.Args().First())
				}
				if err != nil {
					println(err)
					os.Exit(1)
				}
				code, err := pastebin.Put(string(text), c.String("title"))
				if err != nil {
					println("ERROR: ", err.Error())
				}
				println(code)
			},
		},
	}
	app.Run(os.Args)
}
