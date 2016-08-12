package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"yin.mno.stratus.com/gogs/dbulkow/kv"

	"github.com/spf13/cobra"
)

const MacAddrs = "macaddrs/"

var (
	keystore kv.KV
	sortby   string
)

func init() {
	macaddrCmd := &cobra.Command{
		Use:    "macaddr",
		Short:  "List or modify MAC address map",
		Hidden: true,
	}

	listCmd := &cobra.Command{
		Use:     "list [<hostname>]",
		Aliases: []string{"ls"},
		Short:   "List MAC address map",
		Run:     macaddrList,
	}

	setCmd := &cobra.Command{
		Use:   "set <macaddr> <hostname>",
		Short: "Set/Modify MAC a single address mapping",
		Run:   macaddrSet,
	}

	delCmd := &cobra.Command{
		Use:     "delete <macaddr>",
		Aliases: []string{"rm", "del", "remove"},
		Short:   "Remove MAC a single address mapping",
		Run:     macaddrDel,
	}

	saveCmd := &cobra.Command{
		Use:   "save <filename>",
		Short: "Save macaddr map to a file",
		Run:   macaddrSave,
	}

	restoreCmd := &cobra.Command{
		Use:   "restore <filename>",
		Short: "Restore macaddr map from a file",
		Run:   macaddrRestore,
	}

	listCmd.Flags().StringVar(&sortby, "sort-by", "", "Sort by [mac, host]")

	macaddrCmd.AddCommand(listCmd)
	macaddrCmd.AddCommand(setCmd)
	macaddrCmd.AddCommand(delCmd)
	macaddrCmd.AddCommand(saveCmd)
	macaddrCmd.AddCommand(restoreCmd)

	RootCmd.AddCommand(macaddrCmd)

	keystore = &kv.Consul{}
}

func macaddrDel(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	macaddr := args[0]

	if err := keystore.Del(MacAddrs + macaddr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func macaddrSet(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	macaddr := args[0]
	hostname := args[1]

	if err := keystore.Set(MacAddrs+macaddr, hostname); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

type byKey []*kv.KVPair

func (b byKey) Len() int           { return len(b) }
func (b byKey) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byKey) Less(i, j int) bool { return strings.Compare(b[i].Key, b[j].Key) < 0 }

type byVal []*kv.KVPair

func (b byVal) Len() int           { return len(b) }
func (b byVal) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byVal) Less(i, j int) bool { return strings.Compare(b[i].Val, b[j].Val) < 0 }

func macaddrList(cmd *cobra.Command, args []string) {
	filter := ""
	if len(args) == 1 {
		filter = args[0]
	}

	pairs, err := keystore.List(strings.TrimSuffix(MacAddrs, "/"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch sortby {
	case "mac":
		sort.Sort(byKey(pairs))
	case "host":
		sort.Sort(byVal(pairs))
	case "":
	default:
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	for _, p := range pairs {
		if filter != "" && !Glob(filter, p.Val) {
			continue
		}

		fmt.Println(strings.TrimPrefix(p.Key, MacAddrs), p.Val)
	}
}

func macaddrSave(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	filename := args[0]

	pairs, err := keystore.List(MacAddrs)
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

	for _, p := range pairs {
		fmt.Fprintf(file, "%s %s\n", strings.TrimPrefix(p.Key, MacAddrs), p.Val)
	}
}

func macaddrRestore(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	mapfile := args[0]

	file, err := os.Open(mapfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	macaddrs := make(map[string]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		macaddrs[words[0]] = words[1]
	}

	for mac, hostname := range macaddrs {
		if err := keystore.Set(MacAddrs+mac, hostname); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
