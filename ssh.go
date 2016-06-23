package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
)

const SshUsage = `
Usage: lab ssh [OPTIONS] <ftServer> [command]

Exec ssh for an ftServer, when an address is available

Options:
`

func (s *state) ssh(args []string) {
	flagset := flag.NewFlagSet("ssh", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, SshUsage)
		flagset.PrintDefaults()
	}

	sshopt := flagset.String("ssh", "", "ssh command line options")

	if err := flagset.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "flag parse error:", err)
		return
	}

	if flagset.NArg() < 1 {
		flagset.Usage()
		return
	}

	target := strings.Split(flagset.Arg(0), "@")

	u, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to determine username:", err)
		return
	}

	user := u.Username
	mach := ""

	if len(target) == 2 {
		user = target[0]
		mach = target[1]
	} else {
		mach = target[0]
	}

	addr, err := s.getHost(mach)
	if err != nil {
		fmt.Fprintln(os.Stderr, "getHost:", err)
		return
	}

	sshargs := make([]string, 0)
	sshargs = append(sshargs, "ssh")
	if *sshopt != "" {
		sshargs = append(sshargs, strings.Fields(*sshopt)...)
	}
	sshargs = append(sshargs, fmt.Sprintf("%s@%s", user, addr))

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
