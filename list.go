package main

import (
	"flag"
	"fmt"
	"os"
)

const ListUsage = `
Usage: lab list [OPTIONS] [machine filter]

List ftServers
`

func (s *state) list(args []string) {
	flagset := flag.NewFlagSet("list", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, ListUsage)
		flagset.PrintDefaults()
	}

	if err := flagset.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "flag parse error:", err)
		return
	}

	filter := ""
	if flagset.NArg() == 1 {
		filter = flagset.Arg(0)
	}

	reply, err := s.getData(s.labmap, Machines)
	if err != nil {
		fmt.Fprintln(os.Stderr, "labmap(machines):", err)
		return
	}

	machines := reply.Machines

	fmtstr := "%-8s %-3s %-3s %-4s %-5s %-5s\n"

	fmt.Printf(fmtstr, "machine", "cab", "pos", "plat", "power", "state")

	for _, m := range machines {
		if filter != "" && !Glob(filter, m) {
			continue
		}

		reply, err := s.getData(s.labmap, Cabinet+"?machine="+m)
		if err != nil {
			fmt.Fprintln(os.Stderr, "labmap(cabinet):", err)
			return
		}

		cab := reply.Cabinets[m]

		reply, err = s.getData(s.platformid, PlatformID+m)
		if err != nil {
			fmt.Fprintln(os.Stderr, "platformid:", err)
			return
		}

		bmc := reply.BMC

		primary := 0
		if bmc[0].Primary {
			primary = 0
		} else if bmc[1].Primary {
			primary = 1
		} else {
			fmt.Fprintln(os.Stderr, "no primary BMC")
			return
		}

		power := ""
		if bmc[primary].PowerOn == false {
			power = "off"
		}

		fmt.Printf(fmtstr, m, cab.Cabinet, cab.Position, bmc[primary].Platform, power, bmc[primary].State)
	}
}
