package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Machines   = "/v1/machines/"
	Cabinet    = "/v1/cabinet/"
	Address    = "/v1/address/"
	PlatformID = "/v1/platformid/"
)

type cabinet struct {
	VTM0     string `json:"vtm0"`
	VTM1     string `json:"vtm1"`
	Cabinet  string `json:"cabinet"`
	Position string `json:"position"`
	COM1     string `json:"com1"`
	COM2     string `json:"com2"`
	Outlet   string `json:"outlet"`
	KVM      string `json:"kvm"`
	PDU0     string `json:"pdu0"`
	PDU1     string `json:"pdu1"`
}

type address struct {
	MAC     string `json:"macaddr"`
	IP      string `json:"ip"`
	Updated int64  `json:"updated"`
}

type firmware struct {
	BMCSize    int    `json:"bmcsize"`
	BMCA       string `json:"bmca"`
	BMCB       string `json:"bmcb"`
	BootLoader string `json:"bootloader"`
	Running    string `json:"running"`
	SDR        string `json:"sdr"`
	BIOS       string `json:"bios"`
}

type bmc struct {
	Status   string   `json:"status"`
	Error    string   `json:"error,omitempty"`
	Platform string   `json:"platform"`
	PowerOn  bool     `json:"poweron"`
	State    string   `json:"state"`
	CRU      int      `json:"cru"`
	Primary  bool     `json:"primary"`
	Firmware firmware `json:"firmware"`
}

type addr struct {
	MAC     string `json:"macaddr"`
	IP      string `json:"ip"`
	Updated int64  `json:"updated"`
}

type Reply struct {
	Status   string             `json:"status"`
	Error    string             `json:"error"`
	Cabinets map[string]cabinet `json:"cabinets"` /* labmap */
	Machines []string           `json:"machines"`
	Addrs    map[string]address `json:"addrs"`
	Machine  string             `json:"machine"` /* platformid */
	BMC      []*bmc             `json:"bmc"`
}

type MacMapReply struct {
	Status   string            `json:"status"`
	Error    string            `json:"error"`
	Macaddrs map[string]string `json:"macaddrs"`
	Address  map[string]*addr  `json:"addrs"`
}

func (s *state) getData(url, uri string) (*Reply, error) {
	if s.debug {
		fmt.Println("getData", url, uri)
	}

	client := &http.Client{Timeout: time.Second * 20}

	resp, err := client.Get(url + uri)
	if err != nil {
		return nil, fmt.Errorf("connection to %s failed: %v", url+uri, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("service returned code: %s", http.StatusText(resp.StatusCode))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	reply := &Reply{}

	if err := json.Unmarshal(b, reply); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	if reply.Status == "Failed" {
		return nil, fmt.Errorf("request failed: %s", reply.Error)
	}

	return reply, nil
}

func (s *state) getAddr(vtm string) (string, error) {
	url := s.macmap + Address + vtm

	if s.debug {
		fmt.Println("getAddr", url)
	}

	client := &http.Client{Timeout: time.Second * 20}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("connection to %s failed: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("service returned code: %s", http.StatusText(resp.StatusCode))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}

	reply := &MacMapReply{}

	if err := json.Unmarshal(b, reply); err != nil {
		return "", fmt.Errorf("unmarshal: %v", err)
	}

	if reply.Status == "Failed" {
		return "", fmt.Errorf("request failed: %s", reply.Error)
	}

	return reply.Address[vtm].IP, nil
}
