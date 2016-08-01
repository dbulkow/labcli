package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)

const XDGOpen = "xdg-open"

var secondary bool

func init() {
	vtmCmd := &cobra.Command{
		Use:   "vtm <machine>",
		Short: "Start browser window to BMC VTM",
		Run:   vtm,
	}

	vtmCmd.Flags().BoolVar(&secondary, "secondary", false, "Connect to secondar BMC")

	RootCmd.AddCommand(vtmCmd)
}

func vtm(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		return
	}

	mach := args[0]

	reply, err := getData(cmd, "labmap", Cabinet)
	if err != nil {
		fmt.Fprintln(os.Stderr, "labmap(cabinet):", err)
		return
	}

	cab := reply.Cabinets

	vtm := ""
	if secondary {
		vtm = cab[mach].VTM1
	} else {
		vtm = cab[mach].VTM0
	}

	addr, err := getAddr(cmd, vtm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "macmap(address):", err)
		return
	}

	env := os.Environ()

	xdg, err := exec.LookPath(XDGOpen)
	if err != nil {
		fmt.Fprintln(os.Stderr, "xdg-open command not found in PATH")
		return
	}

	xdgargs := make([]string, 0)
	xdgargs = append(xdgargs, XDGOpen)
	xdgargs = append(xdgargs, "http://"+addr)

	err = syscall.Exec(xdg, xdgargs, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s: %v\n", XDGOpen, err)
		return
	}
}
