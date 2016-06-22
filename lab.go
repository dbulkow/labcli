package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	MACMap = "http://yin.mno.stratus.com/"
	LabMap = "http://yin.mno.stratus.com/"
	Etcd   = "http://yin.mno.stratus.com:2379/"
)

type cabinet struct {
	VTM0     string `json:"vtm0"`
	VTM1     string `json:"vtm1"`
	Cabinet  string `json:"cabinet"`
	Position string `json:"position"`
	COM1     string `json:"com1"`
	COM2     string `json:"com2"`
	Outlet   string `json:"outlet"`
	KVM      string `json:"kvm"`
	PDU0     string `json:"pdu0"`
	PDU1     string `json:"pdu1"`
}

type address struct {
	MAC     string `json:"macaddr"`
	IP      string `json:"ip"`
	Updated int64  `json:"updated"`
}

type Reply struct {
	Status   string             `json:"status"`
	Error    string             `json:"error"`
	Cabinets map[string]cabinet `json:"cabinets"`
	Machines []string           `json:"machines"`
	Addrs    map[string]address `json:"addrs"`
}

const Usage = `Usage: lab [OPTIONS] COMMAND [OPTIONS] [arg...]
       lab [ --help | -v | --version ]

Query lab resources.

Options:
`

const CmdUsage = `Commands:
    config   Update/view configuration
    firmware Query firmware versions
    kvm      Connect to KVM
    list     List ftServers
    power    Control CRU power
    ssh      Connect to ftServer using ssh
    telnet   Connect to serial console
    vtm      Start web page to Primary BMC
`

type state struct {
	macmap string
	labmap string
	etcd   string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, Usage)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s\n", CmdUsage)
	}

	//	ConfFile := UserHomeDir() + "/.config/labcli.conf"

	var (
		macmap = flag.String("macmap", MACMap, "URL for macmap")
		labmap = flag.String("labmap", LabMap, "URL for labmap")
		etcd   = flag.String("etcd", Etcd, "URL for etcd")
	)

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	args := flag.Args()

	s := &state{
		macmap: *macmap,
		labmap: *labmap,
		etcd:   *etcd,
	}

	subcmds := map[string]func([]string){
		"cfg":      s.config,
		"config":   s.config,
		"firmware": s.firmware,
		"kvm":      s.kvm,
		"ls":       s.list,
		"list":     s.list,
		"power":    s.power,
		"ssh":      s.ssh,
		"telnet":   s.telnet,
		"vtm":      s.vtm,
	}

	cmd, ok := subcmds[args[0]]
	if !ok {
		flag.Usage()
		return
	}

	cmd(args[1:])
}
