package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh <ftServer> [command]",
	Short: "Exec ssh for an ftServer, when an address is available",
	Run:   ssh,
}

var sshopt string

func init() {
	sshCmd.Flags().StringVar(&sshopt, "opt", "", "ssh command line options")
}

func ssh(cmd *cobra.Command, args []string) {

	target := strings.Split(args[0], "@")

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

	addr, err := getHost(cmd.Flag("etcd").Value.String(), mach)
	if err != nil {
		fmt.Fprintln(os.Stderr, "getHost:", err)
		return
	}

	sshargs := make([]string, 0)
	sshargs = append(sshargs, "ssh")
	if sshopt != "" {
		sshargs = append(sshargs, strings.Fields(sshopt)...)
	}
	sshargs = append(sshargs, fmt.Sprintf("%s@%s", user, addr))

	n := len(args) - 1
	for i := 1; n > 0; i, n = i+1, n-1 {
		sshargs = append(sshargs, args[i])
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
