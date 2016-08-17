package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"

	macapi "yin.mno.stratus.com/gogs/dbulkow/macmap/api"
)

func init() {
	setengCmd := &cobra.Command{
		Use:   "seteng <ftServer-vtm# or IP address>",
		Short: "Enable eng/eng on a BMC",
		Long:  "Enable the eng account (password eng) on a BMC",
		Run:   seteng,
	}

	var sshdebug bool
	setengCmd.Flags().BoolVar(&sshdebug, "debug", false, "Show interactions with server")
	setengCmd.Flags().MarkHidden("debug")

	RootCmd.AddCommand(setengCmd)
}

func setengMapLookup(hostname string) (net.IP, error) {
	addr, err := macapi.GetAddress(macmap, hostname)
	if err != nil {
		return nil, err
	}

	if addr == nil {
		return nil, errors.New("BMC name unknown")
	}

	ip := net.ParseIP(addr.IP)

	if ip == nil {
		return nil, errors.New("IP parse failed")
	}

	return ip, nil
}

func setengNameLookup(hostname string) (net.IP, error) {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(addrs[0])

	if ip == nil {
		return nil, errors.New("IP parse failed")
	}

	return ip, nil
}

func seteng(cmd *cobra.Command, args []string) {
	sshdebug := cmd.Flag("debug").Value.String() == "true"

	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	ip := net.ParseIP(args[0])
	if ip == nil {
		var err error
		ip, err = setengMapLookup(args[0])
		if err != nil {
			maperr := err
			ip, err = setengNameLookup(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Hostname not found:\nmacmap lookup error \"%v\"\nDNS lookup \"%v\"\n", maperr, err)
				os.Exit(1)
			}
		}
	}

	config := &ssh.ClientConfig{
		User: "ADMIN",
		Auth: []ssh.AuthMethod{
			ssh.Password("ADMIN"),
		},
	}

	conn, err := ssh.Dial("tcp", ip.String()+":22", config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 57600, // input speed = 57.6k
		ssh.TTY_OP_OSPEED: 57600, // output speed = 57.6k
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Fprintf(os.Stderr, "request for pty failed: %v\n", err)
		os.Exit(1)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := session.Shell(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	expect := []struct {
		get  string
		send string
	}{
		{"\n-> ", "cd admin1/sp1\r"},
		{"\n-> ", "oemnec ct 20 18 47 03 01\r"},
		{"\n-> ", "exit\r"},
	}

	stripctl := func(str string) string {
		return strings.Map(func(r rune) rune {
			if strings.IndexRune("\r\n", r) < 0 {
				return r
			}
			return -1
		}, str)
	}

	buf := make([]byte, 512)
	end := 0
	token := 0

	for {
		n, err := stdout.Read(buf[end:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			break
		}

		end += n

		if sshdebug {
			fmt.Printf("read %d bytes, %d bytes so far\n%s", n, end, hex.Dump(buf[:end]))
		}

		if bytes.Contains(buf[:end], []byte(expect[token].get)) {
			if sshdebug {
				fmt.Printf("found token %d \"%s\"\n", token, stripctl(expect[token].get))
				fmt.Printf("sending \"%s\"\n", stripctl(expect[token].send))
			}
			fmt.Fprintf(stdin, expect[token].send)
			token++
			end = 0
		}

		if token >= len(expect) {
			break
		}

		if n == 0 {
			time.Sleep(20 * time.Millisecond)
		}
	}
}
