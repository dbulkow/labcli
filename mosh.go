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

var moshopt string

func init() {
	moshCmd := &cobra.Command{
		Use:   "mosh <ftServer> [command]",
		Short: "Exec mosh for an ftServer, when an address is available",
		Run:   moshcmd,
	}

	moshCmd.Flags().StringVar(&moshopt, "opt", "", "mosh command line options")

	RootCmd.AddCommand(moshCmd)
}

func moshcmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		return
	}

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

	addr, err := getHost(cmd, mach)
	if err != nil {
		fmt.Fprintln(os.Stderr, "getHost:", err)
		return
	}

	moshargs := make([]string, 0)
	moshargs = append(moshargs, "mosh")
	if moshopt != "" {
		moshargs = append(moshargs, strings.Fields(moshopt)...)
	}
	moshargs = append(moshargs, fmt.Sprintf("%s@%s", user, addr))

	n := len(args) - 1
	for i := 1; n > 0; i, n = i+1, n-1 {
		moshargs = append(moshargs, args[i])
	}

	env := os.Environ()

	mosh, err := exec.LookPath("mosh")
	if err != nil {
		fmt.Fprintln(os.Stderr, "mosh command not found in PATH")
		return
	}

	err = syscall.Exec(mosh, moshargs, env)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running mosh:", err)
		return
	}
}
