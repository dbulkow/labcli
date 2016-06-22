package main

import "os"

func UserHomeDir() string {
	return os.Getenv("HOME")
}
