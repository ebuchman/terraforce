package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const ValSetAnon = "anon"

var (
	machFlag = cli.StringFlag{
		Name:  "machines",
		Value: "mach[0-3]",
		Usage: "Comma separated list of machine names",
	}

	interactiveFlag = cli.BoolFlag{
		Name:  "interactive,i",
		Usage: "Interactive ssh session",
	}

	iterativeFlag = cli.BoolFlag{
		Name:  "iterative",
		Usage: "Replace '?' in the source path with the machine number",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "terraforce"
	app.Usage = "terraforce [command] [args...]"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{

		{
			Name:  "ssh",
			Usage: "Execute a command through ssh on all machines",
			Flags: []cli.Flag{
				machFlag,
				interactiveFlag,
			},
			Action: func(c *cli.Context) error {
				cmdSsh(c)
				return nil
			},
		},

		{
			Name:  "scp",
			Usage: "Copy a file through scp on all machines",
			Flags: []cli.Flag{
				machFlag,
				iterativeFlag,
			},
			Action: func(c *cli.Context) error {
				cmdScp(c)
				return nil
			},
		},
	}
	app.Run(os.Args)

}

//--------------------------------------------------------------------------------
