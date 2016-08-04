package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	MACMap        = "http://yin.mno.stratus.com/"
	LabMap        = "http://yin.mno.stratus.com/"
	Hosts         = "http://yin.mno.stratus.com/"
	PlatformIDurl = "http://yin.mno.stratus.com/"
)

var RootCmd = &cobra.Command{
	Use:   "lab",
	Short: "Interact with ftLinux lab environment",
}

var (
	macmap     string
	labmap     string
	hosts      string
	platformid string

	verbose bool
)

func main() {
	RootCmd.PersistentFlags().StringVar(&macmap, "macmap", MACMap, "URL for macmap")
	RootCmd.PersistentFlags().StringVar(&labmap, "labmap", LabMap, "URL for labmap")
	RootCmd.PersistentFlags().StringVar(&hosts, "hosts", Hosts, "URL for host lookup")
	RootCmd.PersistentFlags().StringVar(&platformid, "platformid", PlatformIDurl, "URL for platformid")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable communication debugging")

	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
