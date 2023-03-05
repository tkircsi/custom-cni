package main

import (
	"encoding/json"
	"fmt"
	"syscall"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/vishvananda/netlink"
)

type SimpleBridge struct {
	BridgeName string `json:"bridgeName"`
	IP         string `json:"ip"`
}

func main() {
	skel.PluginMain(add, check, del, version.All, "about custom-cni plugin")
}

func add(args *skel.CmdArgs) error {
	sb := SimpleBridge{}
	if err := json.Unmarshal(args.StdinData, &sb); err != nil {
		return err
	}
	fmt.Printf("%v\n", sb)

	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name:   sb.BridgeName,
			MTU:    1500,
			TxQLen: -1,
		},
	}

	err := netlink.LinkAdd(br)
	if err != nil && err != syscall.EEXIST {
		return err
	}

	if err := netlink.LinkSetUp(br); err != nil {
		return err
	}
	return nil
}

func check(args *skel.CmdArgs) error {
	return nil
}

func del(args *skel.CmdArgs) error {
	return nil
}
