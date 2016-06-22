package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	MACMap        = "http://yin.mno.stratus.com/"
	LabMap        = "http://yin.mno.stratus.com/"
	Etcd          = "http://yin.mno.stratus.com:2379/"
	PlatformIDurl = "http://yin.mno.stratus.com/"
)

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
	macmap     string
	labmap     string
	etcd       string
	platformid string

	debug bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, Usage)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s\n", CmdUsage)
	}

	//	ConfFile := UserHomeDir() + "/.config/labcli.conf"

	var (
		macmap     = flag.String("macmap", MACMap, "URL for macmap")
		labmap     = flag.String("labmap", LabMap, "URL for labmap")
		etcd       = flag.String("etcd", Etcd, "URL for etcd")
		platformid = flag.String("platformid", PlatformIDurl, "URL for platformid")
		debug      = flag.Bool("debug", false, "Enable communication debugging")
	)

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	args := flag.Args()

	s := &state{
		macmap:     *macmap,
		labmap:     *labmap,
		etcd:       *etcd,
		platformid: *platformid,
		debug:      *debug,
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
