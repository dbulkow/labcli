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

func main() {
	//	ConfFile := UserHomeDir() + "/.config/labcli.conf"

	var (
		macmap     string
		labmap     string
		etcd       string
		platformid string

		verbose bool
	)

	rootCmd := &cobra.Command{
		Use:   "lab",
		Short: "Query ftServer details",
	}

	rootCmd.PersistentFlags().StringVar(&macmap, "macmap", MACMap, "URL for macmap")
	rootCmd.PersistentFlags().StringVar(&labmap, "labmap", LabMap, "URL for labmap")
	rootCmd.PersistentFlags().StringVar(&etcd, "etcd", Etcd, "URL for etcd")
	rootCmd.PersistentFlags().StringVar(&platformid, "platformid", PlatformIDurl, "URL for platformid")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable communication debugging")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(telnetCmd)
	rootCmd.AddCommand(sshCmd)
	rootCmd.AddCommand(vtmCmd)
	rootCmd.AddCommand(firmwareCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
