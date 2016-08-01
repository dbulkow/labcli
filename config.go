package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

const LabConfig = "labconfig/"

type ComPort struct {
	Enabled  bool   `json:"enabled"`
	Speed    int    `json:"speed,omitempty"`
	Bits     int    `json:"bits,omitempty"`
	StopBits int    `json:"stopbits,omitempty"`
	Parity   string `json:"parity,omitempty"`
	Device   string `json:"device,omitempty"`
}

func (c *ComPort) String() string {
	if !c.Enabled {
		return "no"
	}

	return fmt.Sprintf("%d,%d,%d,%s:%s", c.Speed, c.Bits, c.StopBits, c.Parity, c.Device)
}

type Config struct {
	Name     string  `json:"name"`
	Cabinet  int     `json:"cabinet"`
	Position int     `json:"position"`
	COM1     ComPort `json:"com1"`
	COM2     ComPort `json:"com2"`
	PDU      int     `json:"pdu"`
	KVM      int     `json:"kvm"`
}

var config = &Config{}

func init() {
	configCmd := &cobra.Command{
		Use:    "config",
		Short:  "List or modify lab configuration",
		Hidden: true,
	}

	listCmd := &cobra.Command{
		Use:     "list [<hostname>]",
		Aliases: []string{"ls"},
		Short:   "List lab configuration database",
		Run:     configList,
	}

	setCmd := &cobra.Command{
		Use:   "set <hostname>",
		Short: "Set/Modify lab configuration database",
		Run:   configSet,
	}

	delCmd := &cobra.Command{
		Use:     "delete <hostname>",
		Aliases: []string{"rm", "del", "remove"},
		Short:   "Remove lab config for specified host",
		Run:     configDel,
	}

	saveCmd := &cobra.Command{
		Use:   "save <filename>",
		Short: "Save lab config to a map file",
		Run:   configSave,
	}

	restoreCmd := &cobra.Command{
		Use:   "restore <filename>",
		Short: "Restore lab config from a map file",
		Run:   configRestore,
	}

	setCmd.Flags().IntVar(&config.Cabinet, "cab", 0, "set cabinet number")
	setCmd.Flags().IntVar(&config.Position, "pos", 0, "set position number")
	setCmd.Flags().IntVar(&config.PDU, "pdu", 0, "set PDU slot")
	setCmd.Flags().IntVar(&config.KVM, "kvm", 0, "set KVM slot")
	setCmd.Flags().BoolVar(&config.COM1.Enabled, "com1-enabled", false, "enable COM1")
	setCmd.Flags().IntVar(&config.COM1.Speed, "com1-speed", 57600, "set COM1 port speed")
	setCmd.Flags().IntVar(&config.COM1.Bits, "com1-bits", 8, "set COM1 port bits")
	setCmd.Flags().IntVar(&config.COM1.StopBits, "com1-stop", 1, "set COM1 port stopbits")
	setCmd.Flags().StringVar(&config.COM1.Parity, "com1-parity", "N", "set COM1 port parity [\"N\", \"E\", \"O\"]")
	setCmd.Flags().StringVar(&config.COM1.Device, "com1-device", "", "set COM1 port device on server")
	setCmd.Flags().BoolVar(&config.COM2.Enabled, "com2-enabled", false, "enable COM2")
	setCmd.Flags().IntVar(&config.COM2.Speed, "com2-speed", 57600, "set COM2 port speed")
	setCmd.Flags().IntVar(&config.COM2.Bits, "com2-bits", 8, "set COM2 port bits")
	setCmd.Flags().IntVar(&config.COM2.StopBits, "com2-stop", 1, "set COM2 port stopbits")
	setCmd.Flags().StringVar(&config.COM2.Parity, "com2-parity", "N", "set COM2 port parity [\"N\", \"E\", \"O\"]")
	setCmd.Flags().StringVar(&config.COM2.Device, "com2-device", "", "set COM2 port device on server")

	configCmd.AddCommand(listCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(delCmd)
	configCmd.AddCommand(saveCmd)
	configCmd.AddCommand(restoreCmd)

	RootCmd.AddCommand(configCmd)
}

func configDel(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	hostname := args[0]

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	kv := client.KV()

	key := LabConfig + hostname

	if _, err := kv.Delete(key, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func configSet(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		os.Exit(1)
	}

	hostname := args[0]

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	kv := client.KV()

	key := LabConfig + hostname

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := &Config{}

	if pair != nil {
		if err := json.Unmarshal(pair.Value, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "unmarshal: %v\n", err)
			os.Exit(1)
		}
	}

	cfg.Name = hostname

	if cmd.Flag("cab").Changed {
		cfg.Cabinet = config.Cabinet
	}
	if cmd.Flag("pos").Changed {
		cfg.Position = config.Position
	}
	if cmd.Flag("pdu").Changed {
		cfg.PDU = config.PDU
	}
	if cmd.Flag("kvm").Changed {
		cfg.KVM = config.KVM
	}
	if cmd.Flag("com1-enabled").Changed {
		cfg.COM1.Enabled = config.COM1.Enabled
		cfg.COM1.Speed = config.COM1.Speed
		cfg.COM1.Bits = config.COM1.Bits
		cfg.COM1.StopBits = config.COM1.StopBits
		cfg.COM1.Parity = config.COM1.Parity
		cfg.COM1.Device = config.COM1.Device
	}
	if cmd.Flag("com2-enabled").Changed {
		cfg.COM2.Enabled = config.COM2.Enabled
		cfg.COM2.Speed = config.COM2.Speed
		cfg.COM2.Bits = config.COM2.Bits
		cfg.COM2.StopBits = config.COM2.StopBits
		cfg.COM2.Parity = config.COM2.Parity
		cfg.COM2.Device = config.COM2.Device
	}

	b, err := json.Marshal(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		os.Exit(1)
	}

	_, err = kv.Put(&consul.KVPair{Key: key, Value: b}, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func configList(cmd *cobra.Command, args []string) {
	filter := ""
	if len(args) == 1 {
		filter = args[0]
	}

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	kv := client.KV()

	pairs, _, err := kv.List("labconfig", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg := &Config{}

	header := "%-8s %-4s %-4s %-4s %-4s %-25s %-25s\n"
	format := "%-8s %-4d %-4d %-4d %-4d %-25s %-25s\n"
	fmt.Printf(header, "machine", "cab", "pos", "pdu", "kvm", "com1", "com2")
	for _, p := range pairs {
		if err := json.Unmarshal(p.Value, cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if filter != "" && !Glob(filter, cfg.Name) {
			continue
		}

		fmt.Printf(format, cfg.Name, cfg.Cabinet, cfg.Position, cfg.PDU, cfg.KVM, &cfg.COM1, &cfg.COM2)
	}
}

func configSave(cmd *cobra.Command, args []string) {
	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	kv := client.KV()

	pairs, _, err := kv.List("labconfig", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create file \"%s\": %v\n", filename, err)
		os.Exit(1)
	}
	defer file.Close()

	cfg := &Config{}

	for _, p := range pairs {
		if err := json.Unmarshal(p.Value, cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Fprintf(file, "%-8s lnx%d pos%d com1-%-25s com2-%-25s pdu%d kvm%d\n", cfg.Name, cfg.Cabinet, cfg.Position, &cfg.COM1, &cfg.COM2, cfg.PDU, cfg.KVM)
	}
}

func readMapRow(words []string) (*Config, error) {
	machine := words[0]
	cab, err := strconv.Atoi(strings.TrimPrefix(words[1], "lnx"))
	if err != nil {
		return nil, err
	}
	pos, err := strconv.Atoi(strings.TrimPrefix(words[2], "pos"))
	if err != nil {
		return nil, err
	}
	pdu, err := strconv.Atoi(strings.TrimPrefix(words[5], "pdu"))
	if err != nil {
		return nil, err
	}
	com1 := strings.TrimPrefix(words[3], "com1-")
	com2 := strings.TrimPrefix(words[4], "com2-")
	kvm := 0
	if len(words) == 7 {
		kvm, err = strconv.Atoi(strings.TrimPrefix(words[6], "kvm"))
		if err != nil {
			return nil, err
		}
	}

	c := &Config{
		Name:     machine,
		Cabinet:  cab,
		Position: pos,
		PDU:      pdu,
		KVM:      kvm,
	}

	format := "%d,%d,%d,%1s:%s"

	if com1 != "no" {
		c.COM1 = ComPort{Enabled: true}
		fmt.Sscanf(com1, format, &c.COM1.Speed, &c.COM1.Bits, &c.COM1.StopBits, &c.COM1.Parity, &c.COM1.Device)
	}

	if com2 != "no" {
		c.COM2 = ComPort{Enabled: true}
		fmt.Sscanf(com2, format, &c.COM2.Speed, &c.COM2.Bits, &c.COM2.StopBits, &c.COM2.Parity, &c.COM2.Device)
	}

	return c, nil
}

func readMap(mapfile string) ([]*Config, error) {
	file, err := os.Open(mapfile)
	if err != nil {
		return nil, fmt.Errorf("open map file \"%s\": %v", mapfile, err)
	}
	defer file.Close()

	configs := make([]*Config, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c, err := readMapRow(strings.Fields(scanner.Text()))
		if err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}

	return configs, nil
}

func configRestore(cmd *cobra.Command, args []string) {
	var filename string
	if len(args) > 0 {
		filename = args[0]
	}

	configs, err := readMap(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	kv := client.KV()

	for _, cfg := range configs {
		fmt.Println("restoring", cfg.Name)

		b, err := json.Marshal(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
			os.Exit(1)
		}

		key := LabConfig + cfg.Name

		_, err = kv.Put(&consul.KVPair{Key: key, Value: b}, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
