package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var telnetCmd = &cobra.Command{
	Use:   "telnet <ftServer>",
	Short: "Exec telnet for an ftServer",
	Run:   telnet,
}

var (
	com1  bool
	retry bool
)

func init() {
	telnetCmd.Flags().BoolVarP(&com1, "com1", "1", false, "Use COM1")
	telnetCmd.Flags().BoolVarP(&retry, "retry", "r", false, "Retry connection with telnet server")
}

func telnet(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.UsageFunc()(cmd)
		return
	}

	mach := args[0]

	client := &http.Client{Timeout: time.Second * 20}

	labmap := cmd.Flag("labmap").Value.String()

	resp, err := client.Get(labmap + "/v1/cabinet/?machine=" + mach)
	if err != nil {
		fmt.Fprintln(os.Stderr, "connection to labmap failed:", err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
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
		fmt.Fprintln(os.Stderr, "labmap cabinet request failed:", reply.Error)
		return
	}

	cab := reply.Cabinets[mach]

	cmdline := ""
	if com1 {
		cmdline = cab.COM1
	} else {
		cmdline = cab.COM2
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
