package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const TelnetUsage = `
Usage: lab telnet [OPTIONS] <ftServer>

Exec telnet for an ftServer

Options:
`

func (s *state) telnet(args []string) {
	flagset := flag.NewFlagSet("telnet", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, TelnetUsage)
		flagset.PrintDefaults()
	}

	com1 := flagset.Bool("com1", false, "Use COM1")

	if err := flagset.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "flag parse error:", err)
		return
	}

	if flagset.NArg() < 1 {
		flagset.Usage()
		return
	}

	mach := flagset.Arg(0)

	client := &http.Client{Timeout: time.Second * 20}

	resp, err := client.Get(s.labmap + "/v1/cabinet/?machine=" + mach)
	if err != nil {
		fmt.Fprintln(os.Stderr, "connection to labmap failed:", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read from labmap failed:", err)
		return
	}

	reply := &Reply{}

	if err := json.Unmarshal(b, reply); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal labmap:", err)
		return
	}

	if reply.Status == "Failed" {
		fmt.Fprintln(os.Stderr, "labmap cabinet request failed:", reply.Error)
		return
	}

	cab := reply.Cabinets[mach]

	cmdline := ""
	if *com1 {
		cmdline = cab.COM1
	} else {
		cmdline = cab.COM2
	}

	if cmdline == "" {
		fmt.Fprintln(os.Stderr, "Serial port not configured")
		return
	}

	telargs := strings.Fields(cmdline)
	env := os.Environ()

	telnet, err := exec.LookPath("telnet")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Telnet command not found in PATH")
		return
	}

	err = syscall.Exec(telnet, telargs, env)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running telnet:", err)
		return
	}
}
