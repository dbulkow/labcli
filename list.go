package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const ListUsage = `
Usage: lab list [OPTIONS] [machine filter]

List ftServers

Options:
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

	client := &http.Client{Timeout: time.Second * 20}

	rmach, err := client.Get(s.labmap + "/v1/machines/")
	if err != nil {
		fmt.Fprintln(os.Stderr, "connection to labmap failed:", err)
		return
	}
	defer rmach.Body.Close()

	if rmach.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "labmap returned code:", rmach.StatusCode)
		return
	}

	b, err := ioutil.ReadAll(rmach.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read from labmap failed:", err)
		return
	}

	reply := &Reply{}

	if err := json.Unmarshal(b, reply); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshal labmap:", err)
		return
	}

	if reply.Status == "Failed" {
		fmt.Fprintln(os.Stderr, "labmap machine request failed:", reply.Error)
		return
	}

	for _, m := range reply.Machines {
		rcap, err := client.Get(s.labmap + "/v1/cabinet/?machine=" + m)
		if err != nil {
			fmt.Fprintln(os.Stderr, "connection to labmap failed:", err)
			return
		}
		defer rcap.Body.Close()

		b, err = ioutil.ReadAll(rcap.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "read from labmap failed:", err)
			return
		}

		rpy2 := &Reply{}

		if err := json.Unmarshal(b, rpy2); err != nil {
			fmt.Fprintln(os.Stderr, "unmarshal labmap:", err)
			return
		}

		if rpy2.Status == "Failed" {
			fmt.Fprintln(os.Stderr, "labmap cabinet request failed:", reply.Error)
			return
		}

		cab := rpy2.Cabinets[m]
		fmt.Println(m, cab.Cabinet, cab.Position)
	}
}
