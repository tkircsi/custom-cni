#!/bin/bash

sudo CNI_COMMAND=ADD CNI_CONTAINERID=ns1 \
CNI_NETNS=/var/run/netns/ns1 CNI_IFNAME=eth10 \
CNI_PATH=`pwd` \
./build/custom-cni < config