package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	ipCmd := &cobra.Command{
		Use:   "ip <ftServer>",
		Short: "Exec ip for an ftServer, when an address is available",
		Run:   ip,
	}

	RootCmd.AddCommand(ipCmd)
}

func ip(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		return
	}

	ftserver := args[0]

	addr, err := getHost(cmd.Flag("etcd").Value.String(), ftserver)
	if err != nil {
		fmt.Fprintln(os.Stderr, "getHost:", err)
		return
	}

	fmt.Println(addr)
}
