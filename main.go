package main

import (
	"os"
	"path"

	"github.com/codegangsta/cli"
)

const ValSetAnon = "anon"

var (
	tfOutputFlag = cli.StringFlag{
		Name:  "tfOutput",
		Value: "public_ips",
		Usage: "Name of the terraform output containing list of ips",
	}

	machFlag = cli.StringFlag{
		Name:  "machines",
		Value: "mach[0-3]",
		Usage: "Comma separated list of machine names",
	}

	userFlag = cli.StringFlag{
		Name:  "user",
		Value: "root",
		Usage: "User to ssh in with",
	}

	sshKeyFlag = cli.StringFlag{
		Name:  "ssh-key",
		Value: path.Join(os.Getenv("TF_VAR_key_path"), os.Getenv("TF_VAR_key_name")),
		Usage: "Location of ssh key to use for login",
	}

	interactiveFlag = cli.BoolFlag{
		Name:  "interactive,i",
		Usage: "Interactive ssh session",
	}

	iterativeFlag = cli.BoolFlag{
		Name:  "iterative",
		Usage: "Replace '?' in the source path with the machine number",
	}

	fromFlag = cli.BoolFlag{
		Name:  "from",
		Usage: "Copy from machines",
	}

	recursiveFlag = cli.BoolFlag{
		Name:  "r",
		Usage: "Recursively copy directory",
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
				userFlag,
				sshKeyFlag,
				tfOutputFlag,
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
				userFlag,
				sshKeyFlag,
				tfOutputFlag,
				iterativeFlag,
				fromFlag,
				recursiveFlag,
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
