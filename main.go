package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"syscall"

	"github.com/containernetworking/cni/pkg/skel"
	types100 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

type SimpleBridge struct {
	BridgeName string `json:"bridgeName"`
	IP         string `json:"ip"`
}

func init() {
	// this ensures that main runs only on main thread (thread group leader).
	// since namespace ops (unshare, setns) are done for a single thread, we
	// must ensure that the goroutine does not jump from OS thread to thread
	runtime.LockOSThread()
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
	l, err := netlink.LinkByName(sb.BridgeName)
	if err != nil {
		return fmt.Errorf("could not lookup %q: %v", sb.BridgeName, err)
	}

	newBr, ok := l.(*netlink.Bridge)
	if !ok {
		return fmt.Errorf("%q already exists but is not a bridge", sb.BridgeName)
	}

	// 2. Get the namespace of the container
	netNs, err := ns.GetNS(args.Netns)
	if err != nil {
		return err
	}

	// 3. Create a veth on the container and move the host-end veth to host ns.
	// 4. Attach a host-end veth to linux bridge
	hostIface := &types100.Interface{}
	var handler = func(hostNS ns.NetNS) error {
		vethMac, err := GenerateMac()
		if err != nil {
			return err
		}

		hostVeth, _, err := ip.SetupVeth(args.IfName, 1500, vethMac.String(), hostNS)
		if err != nil {
			return err
		}
		hostIface.Name = hostVeth.Name
		return nil
	}

	if err := netNs.Do(handler); err != nil {
		return err
	}

	hostVeth, err := netlink.LinkByName(hostIface.Name)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetMaster(hostVeth, newBr); err != nil {
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

func GenerateMac() (net.HardwareAddr, error) {
	buf := make([]byte, 6)
	var mac net.HardwareAddr

	_, err := rand.Read(buf)
	if err != nil {
		return net.HardwareAddr{}, err
	}

	// Set the local bit
	buf[0] |= 2

	mac = append(mac, buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])

	return mac, nil
}
