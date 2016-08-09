package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func getHost(cmd *cobra.Command, host string) (string, error) {
	server := cmd.Flag("hosts").Value.String()

	client := &http.Client{Timeout: time.Second * 20}

	resp, err := client.Get(server + "/v1/advertise/?host=" + host)
	if err != nil {
		return "", fmt.Errorf("gethost connect failed: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}

	return string(b), nil
}
