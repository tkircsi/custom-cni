package main

import (
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
)

func main() {
	skel.PluginMain(add, check, del, version.All, "about custom-cni plugin")
}

func add(args *skel.CmdArgs) error {
	return nil
}

func check(args *skel.CmdArgs) error {
	return nil
}

func del(args *skel.CmdArgs) error {
	return nil
}
