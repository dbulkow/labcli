package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [machine filter]",
	Aliases: []string{"ls"},
	Short:   "List ftServers",
	Run:     list,
}

var quiet bool

func init() {
	listCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't display header")
}

func list(cmd *cobra.Command, args []string) {
	filter := ""
	if len(args) == 1 {
		filter = args[0]
	}

	reply, err := getData(cmd, "labmap", Machines)
	if err != nil {
		fmt.Fprintln(os.Stderr, "labmap(machines):", err)
		return
	}

	machines := reply.Machines

	fmtstr := "%-8s %-3s %-3s %-4s %-5s %-5s\n"

	if !quiet {
		fmt.Printf(fmtstr, "machine", "cab", "pos", "plat", "power", "state")
	}

	for _, m := range machines {
		if filter != "" && !Glob(filter, m) {
			continue
		}

		reply, err := getData(cmd, "labmap", Cabinet+"?machine="+m)
		if err != nil {
			fmt.Fprintln(os.Stderr, "labmap(cabinet):", err)
			return
		}

		cab := reply.Cabinets[m]

		var b []*bmc

		if Glob("lin*", m) {
			reply, err = getData(cmd, "platformid", PlatformID+m)
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
