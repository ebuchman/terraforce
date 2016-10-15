package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	. "github.com/tendermint/go-common"
)

// TODO
var SSH_KEY = path.Join(os.Getenv("TF_VAR_key_path"), os.Getenv("TF_VAR_key_name"))
var PORT = 22
var USER = "admin"

//--------------------------------------------------------------------------------

func cmdSsh(c *cli.Context) {
	args := c.Args()
	machines := ResolveMachines(ParseMachines(c.String("machines")))
	if len(machines) == 1 && c.Bool("interactive") {
		sshCmdInteractive(machines[0], args)
	} else {
		cmdBase(args, machines, sshCmd)
	}
}

func cmdScp(c *cli.Context) {
	args := c.Args()
	machines := ResolveMachines(ParseMachines(c.String("machines")))
	fmt.Println(args, machines)
	if c.Bool("iterative") {
		var wg sync.WaitGroup
		for i, mach := range machines {
			wg.Add(1)
			go func(mach string, j int) {
				maybeSleep(len(machines), 2000)
				defer wg.Done()
				scpCmdIterative(mach, args, j)
			}(mach, i)
		}
		wg.Wait()
	} else {
		cmdBase(args, machines, scpCmd)
	}
}

func cmdBase(args []string, machines []string, cmd func(string, []string) error) {
	var wg sync.WaitGroup
	for _, mach := range machines {
		wg.Add(1)
		go func(mach string) {
			maybeSleep(len(machines), 2000)
			defer wg.Done()
			cmd(mach, args)
		}(mach)
	}
	wg.Wait()
}

func sshCmd(mach string, args []string) error {
	args = []string{"-i", SSH_KEY, "-p", Fmt("%d", PORT), Fmt("%s@%s", USER, mach), strings.Join(args, " ")}
	if !runProcess("ssh-cmd-"+mach, "ssh", args, true) {
		return errors.New("Failed to exec ssh command on machine " + mach)
	}
	return nil
}

func sshCmdInteractive(mach string, args []string) error {
	// args = []string{"ssh", mach, strings.Join(args, " ")}
	// any passed args are ignored
	args = []string{"-i", SSH_KEY, "-p", Fmt("%d", PORT), Fmt("%s@%s", USER, mach)}
	_, res := runProcessInteractive("ssh-cmd-"+mach, "ssh", args, true)
	if !res {
		return errors.New("Failed to exec ssh command on machine " + mach)
	}
	return nil
}

func scpCmd(mach string, args []string) error {
	if len(args) != 2 {
		return errors.New("scp expects exactly two args")
	}
	args = []string{"-i", SSH_KEY, "-P", Fmt("%d", PORT), args[0], Fmt("%s@%s:%s", USER, mach, args[1])}
	if !runProcess("ssh-cmd-"+mach, "scp", args, true) {
		return errors.New("Failed to exec scp command on machine " + mach)
	}
	return nil
}

func scpCmdIterative(mach string, args []string, n int) error {
	if len(args) != 2 {
		return errors.New("scp expects exactly two args")
	}
	srcPath := args[0]
	dstPath := args[1]

	// XXX: n+1
	srcPath = strings.Replace(srcPath, "?", Fmt("%d", n+1), -1)

	args = []string{"-i", SSH_KEY, "-P", Fmt("%d", PORT), srcPath, Fmt("%s@%s:%s", USER, mach, dstPath)}
	if !runProcess("ssh-cmd-"+mach, "scp", args, true) {
		return errors.New("Failed to exec scp command on machine " + mach)
	}
	return nil
}

//-----------
