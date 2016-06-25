package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var firmwareCmd = &cobra.Command{
	Use:   "firmware [machine filter]",
	Short: "List firmware versions on ftServers",
	Run:   firmware,
}

func firmware(cmd *cobra.Command, args []string) {
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

	fmtstr := "%-8s %-5s %-5s %-5s %-5s\n"

	fmt.Printf(fmtstr, "machine", "BMC", "Boot", "SDR", "BIOS")

	for _, m := range machines {
		if filter != "" && !Glob(filter, m) {
			continue
		}

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
				Primary: true,
				Firmware: versions{
					Running:    "A",
					BMCA:       "",
					BMCB:       "",
					BootLoader: "",
					SDR:        "",
					BIOS:       "",
				},
			}
			b[1] = &bmc{
				Primary: false,
				Firmware: versions{
					Running:    "A",
					BMCA:       "",
					BMCB:       "",
					BootLoader: "",
					SDR:        "",
					BIOS:       "",
				},
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

		booted := b[primary].Firmware.Running
		bmcver := ""
		switch booted {
		case "A":
			bmcver = b[primary].Firmware.BMCA
		case "B":
			bmcver = b[primary].Firmware.BMCB
		}
		bootldr := b[primary].Firmware.BootLoader
		sdr := b[primary].Firmware.SDR
		bios := b[primary].Firmware.BIOS

		fmt.Printf(fmtstr, m, bmcver, bootldr, sdr, bios)
	}
}
