package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"yin.mno.stratus.com/gogs/dbulkow/macmap/api"
)

const MacMap = "macmap/"

func init() {
	macmapCmd := &cobra.Command{
		Use:    "macmap",
		Short:  "List or modify address map",
		Hidden: true,
	}

	listCmd := &cobra.Command{
		Use:     "list [<hostname>]",
		Aliases: []string{"ls"},
		Short:   "List address map",
		Run:     macmapList,
	}

	setCmd := &cobra.Command{
		Use:   "set <hostname> <IP address>",
		Short: "Set/Modify a single address mapping",
		Run:   macmapSet,
	}

	delCmd := &cobra.Command{
		Use:     "delete <hostname>",
		Aliases: []string{"rm", "del", "remove"},
		Short:   "Remove a single address mapping",
		Run:     macmapDel,
	}

	saveCmd := &cobra.Command{
		Use:   "save <filename>",
		Short: "Save map to a file",
		Run:   macmapSave,
	}

	restoreCmd := &cobra.Command{
		Use:   "restore <filename>",
		Short: "Restore map from a file",
		Run:   macmapRestore,
	}

	macmapCmd.AddCommand(listCmd)
	macmapCmd.AddCommand(setCmd)
	macmapCmd.AddCommand(delCmd)
	macmapCmd.AddCommand(saveCmd)
	macmapCmd.AddCommand(restoreCmd)

	RootCmd.AddCommand(macmapCmd)
}

func macmapDel(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	hostname := args[0]

	if err := keystore.Del(MacMap + hostname); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func macmapSet(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	hostname := args[0]
	ip := args[1]

	addr := &api.Addr{IP: ip}

	b, err := json.Marshal(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}

	if err := keystore.Set(MacMap+hostname, string(b)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func macmapList(cmd *cobra.Command, args []string) {
	filter := ""
	if len(args) == 1 {
		filter = args[0]
	}

	pairs, err := keystore.List(strings.TrimSuffix(MacMap, "/"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, p := range pairs {
		hostname := strings.TrimPrefix(p.Key, MacMap)

		if filter != "" && !Glob(filter, hostname) {
			continue
		}

		var addr api.Addr

		if err := json.Unmarshal([]byte(p.Val), &addr); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("%-16s %-16s\n", hostname, addr.IP)
	}
}

func macmapSave(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	filename := args[0]

	pairs, err := keystore.List(MacMap)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create file \"%s\": %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	macmap := make(map[string]*api.Addr)

	for _, p := range pairs {
		addr := &api.Addr{}

		if err := json.Unmarshal([]byte(p.Val), addr); err != nil {
			fmt.Fprintf(os.Stderr, "unmarshal: %v\n", err)
			os.Exit(1)
		}

		macmap[strings.TrimPrefix(p.Key, MacMap)] = addr
	}

	b, err := json.Marshal(macmap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}

	file.Write(b)
}

func macmapRestore(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	mapfile := args[0]

	macmap := make(map[string]*api.Addr)

	b, err := ioutil.ReadFile(mapfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "readfile: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(b, &macmap); err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal: %v\n", err)
		os.Exit(1)
	}

	for hostname, addr := range macmap {
		b, err := json.Marshal(addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
			os.Exit(1)
		}

		if err := keystore.Set(MacMap+hostname, string(b)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
