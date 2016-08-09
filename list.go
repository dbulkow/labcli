package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	labapi "yin.mno.stratus.com/gogs/dbulkow/labmap/api"
	pid "yin.mno.stratus.com/gogs/dbulkow/platformid/api"
)

var quiet bool

func init() {
	listCmd := &cobra.Command{
		Use:     "list [machine filter]",
		Aliases: []string{"ls"},
		Short:   "List ftServers",
		Run:     list,
	}

	listCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't display header")

	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) {
	filter := ""
	if len(args) == 1 {
		filter = args[0]
	}

	machines, err := labapi.Machines(labmap)
	if err != nil {
		fmt.Fprintln(os.Stderr, "labmap(machines):", err)
		return
	}

	fmtstr := "%-8s %-3s %-3s %-4s %-5s %-10s\n"

	if !quiet {
		fmt.Printf(fmtstr, "machine", "cab", "pos", "plat", "power", "state")
	}

	for _, m := range machines {
		if filter != "" && !Glob(filter, m) {
			continue
		}

		cab, err := labapi.GetCabinet(labmap, m)
		if err != nil {
			fmt.Fprintln(os.Stderr, "labmap(cabinet):", err)
			return
		}

		var b []*pid.BMC

		if Glob("lin*", m) {
			p, err := pid.PlatformID(platformid, m)
			if err != nil {
				b = []*pid.BMC{
					{
						Primary:  true,
						Platform: "",
						PowerOn:  true,
						State:    "",
					},
					{
						Primary:  false,
						Platform: "",
						PowerOn:  true,
						State:    "",
					},
				}
			} else {
				b = p.Bmc
			}
		} else {
			b = make([]*pid.BMC, 2)
			b[0] = &pid.BMC{
				Primary:  true,
				Platform: "",
				PowerOn:  true,
				State:    "",
			}
			b[1] = &pid.BMC{
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
