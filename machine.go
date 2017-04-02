package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	. "github.com/tendermint/go-common"
)

// TODO
var PORT = 22

//--------------------------------------------------------------------------------

func cmdSsh(c *cli.Context) {
	args := c.Args()
	machines := ResolveMachines(ParseMachines(c.String("machines")), c.String("tfOutput"))
	if len(machines) == 1 && c.Bool("interactive") {
		sshCmdInteractive(machines[0], c.String("user"), c.String("ssh-key"), args)
	} else {
		if c.Bool("background") {
			cmdBase(args, machines, c.String("user"), c.String("ssh-key"), sshCmdBg)
		} else {
			cmdBase(args, machines, c.String("user"), c.String("ssh-key"), sshCmd)
		}
	}
}

func cmdScp(c *cli.Context) {
	args := c.Args()
	machines := ResolveMachines(ParseMachines(c.String("machines")), c.String("tfOutput"))
	fmt.Println(args, machines)
	if c.Bool("iterative") {
		var wg sync.WaitGroup
		for i, mach := range machines {
			wg.Add(1)
			go func(mach string, j int) {
				maybeSleep(len(machines), 2000)
				defer wg.Done()
				scpCmdIterative(mach, c.String("user"), c.String("ssh-key"), args, j, c.Bool("from"), c.Bool("r"))
			}(mach, i)
		}
		wg.Wait()
	} else {
		cmdBase(args, machines, c.String("user"), c.String("ssh-key"), scpCmdFunc(c.Bool("from"), c.Bool("r"), c.String("user"), c.String("ssh-key")))
	}
}

func cmdBase(args []string, machines []string, user, sshKey string, cmd func(string, string, string, []string) error) {
	var wg sync.WaitGroup
	for _, mach := range machines {
		wg.Add(1)
		go func(mach string) {
			maybeSleep(len(machines), 2000)
			defer wg.Done()
			cmd(mach, user, sshKey, args)
		}(mach)
	}
	wg.Wait()
}

// TODO XXX
// http://www.daemonology.net/blog/2012-01-16-automatically-populating-ssh-known-hosts.html
var OptDisableHostCheck = "StrictHostKeyChecking=no"
var OptNullHostFile = "UserKnownHostsFile=/dev/null"

func sshArgs(mach, user, sshKey string, scp bool) []string {
	p := "-p"
	if scp {
		p = "-P"
	}
	key_port := []string{"-i", sshKey, p, Fmt("%d", PORT)}
	opts := []string{"-o", OptDisableHostCheck, "-o", OptNullHostFile}
	user = Fmt("%s@%s", user, mach)
	args := append(key_port, opts...)
	if !scp {
		args = append(args, user)
	}
	return args
}

func sshCmd(mach, user, sshKey string, args []string) error {
	args = append(sshArgs(mach, user, sshKey, false), args...)
	if !runProcess("ssh-cmd-"+mach, "ssh", args, true) {
		return errors.New("Failed to exec ssh command on machine " + mach)
	}
	return nil
}

func sshCmdBg(mach, user, sshKey string, args []string) error {
	args = append(sshArgs(mach, user, sshKey, false), args...)
	if !runProcessBg("ssh-cmd-"+mach, "ssh", args, true) {
		return errors.New("Failed to exec ssh command on machine " + mach)
	}
	return nil
}

func sshCmdInteractive(mach, user, sshKey string, args []string) error {
	// args = []string{"ssh", mach, strings.Join(args, " ")}
	// any passed args are ignored
	args = sshArgs(mach, user, sshKey, false)
	_, res := runProcessInteractive("ssh-cmd-"+mach, "ssh", args, true)
	if !res {
		return errors.New("Failed to exec ssh command on machine " + mach)
	}
	return nil
}

func scpCmdFunc(from, recursive bool, user, sshKey string) func(string, string, string, []string) error {
	return func(mach, u, s string, args []string) error {
		return scpCmd(mach, user, sshKey, args, from, recursive)
	}

}

func scpCmd(mach, user, sshKey string, args []string, from, recursive bool) error {
	if len(args) != 2 {
		return errors.New("scp expects exactly two args")
	}
	src, dst := args[0], args[1]
	var cpArgs []string
	cpArgs = cpToFrom(src, dst, user, mach, from, recursive)

	args = append(sshArgs(mach, user, sshKey, true), cpArgs...)
	if !runProcess("ssh-cmd-"+mach, "scp", args, true) {
		return errors.New("Failed to exec scp command on machine " + mach)
	}
	return nil

}

// copy to remotes or from remotes
func cpToFrom(src, dst, user, mach string, from, recursive bool) []string {
	var args []string
	if from {
		args = []string{Fmt("%s@%s:%s", user, mach, src), dst}
	} else {
		args = []string{src, Fmt("%s@%s:%s", user, mach, dst)}
	}

	if recursive {
		args = append([]string{"-r"}, args...)
	}
	return args
}

func scpCmdIterative(mach, user, sshKey string, args []string, n int, from, recursive bool) error {
	if len(args) != 2 {
		return errors.New("scp expects exactly two args")
	}
	src, dst := args[0], args[1]
	if from {
		dst = strings.Replace(dst, "?", Fmt("%d", n), -1)
	} else {
		src = strings.Replace(src, "?", Fmt("%d", n), -1)
	}
	cpArgs := cpToFrom(src, dst, user, mach, from, recursive)
	args = append(sshArgs(mach, user, sshKey, true), cpArgs...)
	if !runProcess("ssh-cmd-"+mach, "scp", args, true) {
		return errors.New("Failed to exec scp command on machine " + mach)
	}
	return nil
}

//-----------
