package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	MACMap        = "http://yin.mno.stratus.com/"
	LabMap        = "http://yin.mno.stratus.com/"
	Etcd          = "http://yin.mno.stratus.com:2379/"
	PlatformIDurl = "http://yin.mno.stratus.com/"
)

var RootCmd = &cobra.Command{
	Use:   "lab",
	Short: "Interact with ftLinux lab environment",
}

var (
	macmap     string
	labmap     string
	etcd       string
	platformid string
	kvstore    string

	verbose bool
)

func main() {
	RootCmd.PersistentFlags().StringVar(&macmap, "macmap", MACMap, "URL for macmap")
	RootCmd.PersistentFlags().StringVar(&labmap, "labmap", LabMap, "URL for labmap")
	RootCmd.PersistentFlags().StringVar(&etcd, "etcd", Etcd, "URL for etcd")
	RootCmd.PersistentFlags().StringVar(&platformid, "platformid", PlatformIDurl, "URL for platformid")
	RootCmd.PersistentFlags().StringVar(&kvstore, "kv", "consul", "Select key-value store from [etcd, consul]")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable communication debugging")

	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
