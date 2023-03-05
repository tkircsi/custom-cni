package main

import (
	"encoding/json"
	"fmt"
	"syscall"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/vishvananda/netlink"
	"github.com/containernetworking/plugins/pkg/ns"
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


    // 1. Prepare the netlink.Bridge object we want.
    // 2. Create the Bridge
    // 3. Setup the Linux Bridge.
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


    // 1. Get the bridge object from the Bridge we created before
    // 2. Get the namespace of the container
    // 3. Create a veth on the container and move the host-end veth to host ns.
    // 4. Attach a host-end veth to linux bridge
	l, err := netlink.LinkByName(sb.BridgeName)
	if err != nil {
		return fmt.Errorf("could not lookup %q: %v", sb.BridgeName, err)
	}

	newBr, ok := l.(*netlink.Bridge)
	if !ok {
		return fmt.Errorf("%q already exists but is not a bridge", sb.BridgeName)
	}

	netNs, err := ns.GetNS(args.Netns)
	if err != nil {
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
