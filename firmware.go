package main

import (
	"flag"
	"fmt"
	"os"
)

const FirmwareUsage = `
Usage: lab firmware [OPTIONS] [machine filter]

List firmware versions on ftServers
`

func (s *state) firmware(args []string) {
	flagset := flag.NewFlagSet("firmware", flag.ExitOnError)

	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, FirmwareUsage)
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

	fmtstr := "%-8s %-5s %-5s %-5s %-5s\n"

	fmt.Printf(fmtstr, "machine", "BMC", "Boot", "SDR", "BIOS")

	for _, m := range machines {
		if filter != "" && !Glob(filter, m) && !Glob("lin*", m) {
			continue
		}

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
				Primary: true,
				Firmware: firmware{
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
				Firmware: firmware{
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
