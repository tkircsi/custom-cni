#!/bin/bash

#!/bin/sh
NS1=dev
BR=v-net-0
CID=dev
CNS=/var/run/netns/dev
CIF=veth-dev

sudo ip netns del $NS1
sudo ifconfig $BR down
sudo brctl delbr $BR
sudo ip netns add $NS1
go build -o build/custom-cni .

echo "Ready to call the step3 example"
sudo CNI_COMMAND=ADD CNI_CONTAINERID=$CID CNI_NETNS=$CNS CNI_IFNAME=$CIF CNI_PATH=`pwd` ./build/custom-cni < config
echo "The CNI has been called, see the following results"
echo "The bridge and the veth has been attatch to"
sudo brctl show $BR
echo "The interface in the netns"
sudo ip netns exec $NS1 ifconfig -a