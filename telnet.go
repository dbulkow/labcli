package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	labapi "yin.mno.stratus.com/gogs/dbulkow/labmap/api"
)

var (
	com1  bool
	retry bool
)

func init() {
	telnetCmd := &cobra.Command{
		Use:   "telnet <ftServer>",
		Short: "Exec telnet for an ftServer",
		Run:   telnet,
	}

	telnetCmd.Flags().BoolVarP(&com1, "com1", "1", false, "Use COM1")
	telnetCmd.Flags().BoolVarP(&retry, "retry", "r", false, "Retry connection with telnet server")

	RootCmd.AddCommand(telnetCmd)
}

func telnet(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		return
	}

	machine := args[0]

	cab, err := labapi.GetCabinet(labmap, machine)
	if err != nil {
		fmt.Fprintln(os.Stderr, "connection to labmap failed:", err)
		return
	}

	/* start with COM2, but if not configured try COM1 */
	cmdline := ""
	if cab.COM2 != "" {
		cmdline = cab.COM2
	} else if cab.COM1 != "" {
		cmdline = cab.COM1
	}

	/* user knows best, let them override */
	if com1 {
		cmdline = cab.COM1
	}

	if cmdline == "" {
		fmt.Fprintln(os.Stderr, "Serial port not configured")
		return
	}

	// XXX replace middle argument (machine name) with IP

	telargs := strings.Fields(cmdline)
	env := os.Environ()

	telnet, err := exec.LookPath("telnet")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Telnet command not found in PATH")
		return
	}

	sigs := make(chan os.Signal, 1)

	registerSignals(sigs)

	go func() {
		for {
			sig := <-sigs
			fmt.Println(sig)
			if sig.String() == "interrupt" {
				fmt.Println("leaving")
				os.Exit(1)
			}
		}
	}()
again:
	cmnd := exec.Command(telnet, telargs[1:]...)
	cmnd.Stderr = os.Stderr
	cmnd.Stdout = os.Stdout
	cmnd.Stdin = os.Stdin
	cmnd.Env = env

	fmt.Println(cmnd.Args)

	err = cmnd.Run()
	if err != nil {
		e, ok := err.(*exec.ExitError)
		if !ok {
			fmt.Fprintln(os.Stderr, "run failed:", err)
			return
		}
		wstat := e.Sys().(syscall.WaitStatus)
		fmt.Printf("telnet exit status: %d\n", wstat.ExitStatus())

		if retry {
			time.Sleep(time.Second * 5)
			goto again
		}
	}
}
