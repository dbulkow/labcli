package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const XDGOpen = "xdg-open"

const VTMUsage = `
Usage:  lab vtm [OPTIONS] <machine>

Start browser window to BMC VTM

Options:
`

func (s *state) vtm(args []string) {
	flagset := flag.NewFlagSet("vtm", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, VTMUsage)
		flagset.PrintDefaults()
	}

	secondary := flagset.Bool("secondary", false, "Connect to secondar BMC")

	if err := flagset.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "flag parse error:", err)
		return
	}

	if flagset.NArg() < 1 {
		flagset.Usage()
		return
	}

	mach := flagset.Arg(0)

	reply, err := s.getData(s.labmap, Cabinet)
	if err != nil {
		fmt.Fprintln(os.Stderr, "labmap(cabinet):", err)
		return
	}

	cab := reply.Cabinets

	vtm := ""
	if *secondary {
		vtm = cab[mach].VTM1
	} else {
		vtm = cab[mach].VTM0
	}

	addr, err := s.getAddr(vtm)
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
