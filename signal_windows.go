// +build windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func registerSignals(sigs chan os.Signal) {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
}
