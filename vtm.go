package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"

	labapi "yin.mno.stratus.com/gogs/dbulkow/labmap/api"
	macapi "yin.mno.stratus.com/gogs/dbulkow/macmap/api"
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

	cab, err := labapi.Cabinets(labmap)
	if err != nil {
		fmt.Fprintln(os.Stderr, "labmap(cabinet):", err)
		return
	}

	vtm := ""
	if c, ok := cab[mach]; ok {
		if secondary {
			vtm = c.VTM1
		} else {
			vtm = c.VTM0
		}
	}

	addr, err := macapi.GetAddress(macmap, vtm)
	if err != nil {
		fmt.Fprintln(os.Stderr, "macmap(address):", err)
		return
	}

	if addr == nil {
		fmt.Fprintf(os.Stderr, "address[%s] not found\n", vtm)
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
	xdgargs = append(xdgargs, "http://"+addr.IP)

	err = syscall.Exec(xdg, xdgargs, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s: %v\n", XDGOpen, err)
		return
	}
}
