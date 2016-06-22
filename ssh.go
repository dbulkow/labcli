package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
	"time"
)

const SshUsage = `
Usage: lab ssh [OPTIONS] <ftServer> [command]

Exec ssh for an ftServer, when an address is available

Options:
`

func (s *state) ssh(args []string) {
	flagset := flag.NewFlagSet("ssh", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, ListUsage)
		flagset.PrintDefaults()
	}

	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to determine username:", err)
		return
	}

	user := flagset.String("user", u.Username, "User name to pass to ssh")
	sshopt := flagset.String("ssh", "", "ssh command line options")

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

	resp, err := client.Get(s.etcd + "/v2/keys/hosts/" + mach)
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

	reply := &struct {
		Action string `json:"action"`
		Node   struct {
			Key           string    `json:"key"`
			Value         string    `json:"value"`
			Expiration    time.Time `json:"expiration"`
			TTL           int64     `json:"ttl"`
			ModifiedIndex int       `json:"modifiedindex"`
			CreatedIndex  int       `json:"createdindex"`
		} `json:"node"`
	}{}

	if err := json.Unmarshal(b, reply); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal labmap:", err)
		return
	}

	if reply.Node.Value == "" {
		fmt.Fprintln(os.Stderr, "ssh address not available")
		return
	}

	sshargs := make([]string, 0)
	sshargs = append(sshargs, "ssh")
	sshargs = append(sshargs, strings.Fields(*sshopt)...)
	sshargs = append(sshargs, fmt.Sprintf("%s@%s", *user, reply.Node.Value))

	n := flagset.NArg() - 1
	for i := 1; n > 0; i, n = i+1, n-1 {
		sshargs = append(sshargs, flagset.Arg(i))
	}

	env := os.Environ()

	ssh, err := exec.LookPath("ssh")
	if err != nil {
		fmt.Fprintln(os.Stderr, "ssh command not found in PATH")
		return
	}

	err = syscall.Exec(ssh, sshargs, env)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running ssh:", err)
		return
	}
}
