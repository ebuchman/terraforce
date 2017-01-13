package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	. "github.com/tendermint/go-common"
	pcm "github.com/tendermint/go-process"
	"github.com/tendermint/go-wire"
)

func runProcess(label string, command string, args []string, verbose bool) bool {
	_, res := runProcessGetResult(label, command, args, verbose)
	return res
}

func runProcessGetResult(label string, command string, args []string, verbose bool) (string, bool) {
	outFile := NewBufferCloser(nil)
	proc, err := pcm.StartProcess(label, "", command, args, nil, outFile)
	if err != nil {
		if verbose {
			fmt.Println(Red(err.Error()))
		}
		return "", false
	}

	<-proc.WaitCh
	if verbose {
		fmt.Println(Green(command), Green(args))
	}
	if proc.ExitState.Success() {
		if verbose {
			fmt.Println(Blue(string(outFile.Bytes())))
		}
		return string(outFile.Bytes()), true
	} else {
		// Error!
		if verbose {
			fmt.Println(Red(string(outFile.Bytes())))
		}
		return string(outFile.Bytes()), false
	}
}

func runProcessInteractive(label string, command string, args []string, verbose bool) (string, bool) {
	proc, err := pcm.StartProcess(label, "", command, args, os.Stdin, os.Stdout)
	if err != nil {
		if verbose {
			fmt.Println(Red(err.Error()))
		}
		return "", false
	}

	<-proc.WaitCh
	if verbose {
		fmt.Println(Green(command), Green(args))
	}
	return "", true
}

//--------------------------------------------------------------------------------

func eB(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `$`, `\$`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	s = strings.Replace(s, `'`, `\'`, -1)
	s = strings.Replace(s, `!`, `\!`, -1)
	s = strings.Replace(s, `#`, `\#`, -1)
	s = strings.Replace(s, `%`, `\%`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	s = strings.Replace(s, "`", "\\`", -1)
	return s
}

func condenseBash(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	lines := strings.Split(cmd, "\n")
	res := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		res = append(res, line)
	}
	return strings.Join(res, "; ")
}

//--------------------------------------------------------------------------------

func ReadJSONFile(o interface{}, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	wire.ReadJSON(o, b, &err)
	if err != nil {
		return err
	}
	return nil
}

//--------------------------------------------------------------------------------
// for avoiding request limits in concurrent calls to amazonec2

func maybeSleep(n int, t int) {
	if n > 16 {
		sleepMS := int(rand.Int31n(int32(t)))
		time.Sleep(time.Millisecond * time.Duration(sleepMS))
	}
}

//------------------------------

// read ips out of terraform output
func terraformGetVar(varName string) ([]string, error) {
	args := []string{"output", varName}
	label := "terraformDNS-cmd"
	cmd := "terraform"
	outFile := NewBufferCloser(nil)
	proc, err := pcm.StartProcess(label, "", cmd, args, nil, outFile)
	if err != nil {
		return nil, err
	}

	var result string
	<-proc.WaitCh
	if proc.ExitState.Success() {
		result = string(outFile.Bytes())
	} else {
		// Error!
		return nil, errors.New(string(outFile.Bytes()))
	}

	result = strings.Replace(result, "\n", "", -1)
	spl := strings.Split(result, ",")
	return spl, nil
}
