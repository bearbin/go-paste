package main

import (
	"errors"
	"github.com/bearbin/go-paste/fpaste"
	"github.com/bearbin/go-paste/pastebin"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
)

var errUnknownService = errors.New("unknown paste service")

func main() {
	app := cli.NewApp()
	app.Usage = "get and put pastes from pastebin and other paste sites."
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "service, s", Value: "pastebin", Usage: "the pastebin service to use"},
	}
	app.Commands = []cli.Command{
		{
			Name:      "put",
			ShortName: "p",
			Usage:     "put a paste",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "id", Usage: "return the paste id not the url"},
				cli.StringFlag{Name: "title, t", Value: "", Usage: "the title for the paste"},
			},
			Action: func(c *cli.Context) {
				srv, url, err := convertService(c.GlobalString("service"))
				if err != nil {
					println("ERROR:", err.Error())
					os.Exit(1)
				}
				var text []byte
				if c.Args().First() == "-" || c.Args().First() == "" {
					text, err = ioutil.ReadAll(os.Stdin)
				} else {
					text, err = ioutil.ReadFile(c.Args().First())
				}
				if err != nil {
					println("ERROR:", err.Error())
					os.Exit(1)
				}
				code, err := srv.Put(url, string(text), c.String("title"))
				if err != nil {
					println("ERROR:", err.Error())
					os.Exit(1)
				}
				if c.Bool("id") {
					println(code)
				} else {
					println(srv.WrapID(code))
				}
			},
		},
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "get a paste from its url",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "id", Usage: "get a paste from its ID instead of its URL"},
			},
			Action: func(c *cli.Context) {
				srv, url, err := convertService(c.GlobalString("service"))
				if err != nil {
					println("ERROR:", err.Error())
					os.Exit(1)
				}
				var id string
				if c.Bool("id") {
					id = c.Args().First()
				} else {
					id = srv.StripURL(c.Args().First())
				}
				text, err := srv.Get(url, id)
				if err != nil {
					println("ERROR:", err.Error())
					os.Exit(1)
				}
				println(text)
			},
		},
	}
	app.Run(os.Args)
}

func convertService(srv string) (service, string, error) {
	switch {
	case srv == "pastebin" || srv == "pastebin.com" || srv == "http://pastebin.com":
		return pastebin.Pastebin{}, "", nil
	case srv == "fpaste" || srv == "fpaste.org" || srv == "http://fpaste.org":
		return fpaste.Fpaste{}, "", nil
	}
	return nil, "", errUnknownService
}
