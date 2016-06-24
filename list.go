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

	quiet := flagset.Bool("q", false, "Don't display header")

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

	if !*quiet {
		fmt.Printf(fmtstr, "machine", "cab", "pos", "plat", "power", "state")
	}

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

		var b []*bmc

		if Glob("lin*", m) {
			reply, err = s.getData(s.platformid, PlatformID+m)
			if err != nil {
				fmt.Fprintln(os.Stderr, "platformid:", err)
				return
			}

			b = reply.BMC
		} else {
			b = make([]*bmc, 2)
			b[0] = &bmc{
				Primary:  true,
				Platform: "",
				PowerOn:  true,
				State:    "",
			}
			b[1] = &bmc{
				Primary:  false,
				Platform: "",
				PowerOn:  true,
				State:    "",
			}
		}

		primary := 0
		if b[0].Primary {
			primary = 0
		} else if b[1].Primary {
			primary = 1
		} else {
			fmt.Fprintln(os.Stderr, "no primary BMC")
			return
		}

		power := ""
		if b[primary].PowerOn == false {
			power = "off"
		}

		fmt.Printf(fmtstr, m, cab.Cabinet, cab.Position, b[primary].Platform, power, b[primary].State)
	}
}
