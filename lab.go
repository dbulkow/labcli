package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	BaseURL       = "http://yin.mno.stratus.com"
	MACMap        = BaseURL
	LabMap        = BaseURL
	Hosts         = BaseURL
	PlatformIDurl = BaseURL
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

	RootCmd.PersistentFlags().MarkHidden("macmap")
	RootCmd.PersistentFlags().MarkHidden("labmap")
	RootCmd.PersistentFlags().MarkHidden("hosts")
	RootCmd.PersistentFlags().MarkHidden("platformid")

	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
