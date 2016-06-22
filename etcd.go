package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (s *state) getHost(mach string) (string, error) {
	client := &http.Client{Timeout: time.Second * 20}

	resp, err := client.Get(s.etcd + "/v2/keys/hosts/" + mach)
	if err != nil {
		return "", fmt.Errorf("etcd connect failed: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}

	reply := &struct {
		Action string `json:"action"`
		Node   struct {
			Key           string    `json:"key"`
			Value         string    `json:"value"`
			Expiration    time.Time `json:"expiration"`
			TTL           int64     `json:"ttl"`
			ModifiedIndex int       `json:"modifiedindex"`
			CreatedIndex  int       `json:"createdindex"`
		} `json:"node"`
	}{}

	if err := json.Unmarshal(b, reply); err != nil {
		return "", fmt.Errorf("unmarshal: %v", err)
	}

	if reply.Node.Value == "" {
		return "", fmt.Errorf("ssh address not available")
	}

	return reply.Node.Value, nil
}
