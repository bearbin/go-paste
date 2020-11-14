package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/bearbin/go-paste/fpaste"
	"github.com/bearbin/go-paste/pastebin"
)

var errUnknownService = errors.New("unknown paste service")

func main() {
	app := cli.NewApp()
	app.Usage = "get and put pastes from pastebin and other paste sites."
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "service, s", Value: "pastebin", Usage: "the pastebin service to use"},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "put",
			Usage:   "put a paste",
			Aliases: []string{"p"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "id", Usage: "return the paste id not the url"},
				&cli.StringFlag{Name: "title, t", Value: "", Usage: "the title for the paste"},
			},
			Action: func(c *cli.Context) error {
				srv, err := convertService(c.String("service"))
				if err != nil {
					return err
				}

				var text []byte
				if c.Args().First() == "-" || c.Args().First() == "" {
					text, err = ioutil.ReadAll(os.Stdin)
				} else {
					text, err = ioutil.ReadFile(c.Args().First())
				}

				if err != nil {
					return err
				}

				code, err := srv.Put(string(text), c.String("title"))
				if err != nil {
					return err
				}

				if !c.Bool("id") {
					code = srv.WrapID(code)
				}

				_, _ = fmt.Fprintln(app.Writer, code)

				return nil
			},
		},
		{
			Name:    "get",
			Usage:   "get a paste from its url",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "id", Usage: "get a paste from its ID instead of its URL"},
			},
			Action: func(c *cli.Context) error {
				srv, err := convertService(c.String("service"))
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
				text, err := srv.Get(id)
				if err != nil {
					return err
				}

				_, _ = fmt.Fprintln(app.Writer, text)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(app.ErrWriter, "%v\n", err)

		os.Exit(1)
	}
}

func convertService(srv string) (service, error) {
	switch {
	case srv == "pastebin" || srv == "pastebin.com" || srv == "http://pastebin.com":
		return pastebin.Pastebin{}, nil
	case srv == "fpaste" || srv == "fpaste.org" || srv == "http://fpaste.org":
		return fpaste.Fpaste{}, nil
	}
	return nil, errUnknownService
}
