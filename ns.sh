#!/bin/bash

# Demonstrates network namespaces and visibility.
set euxo -pipefail
ip netns add dev
ip netns add prod
ip link add v-net-0 type bridge
ip link set dev v-net-0 up
ip link add veth-dev type veth peer name veth-dev-br
ip link add veth-prod type veth peer name veth-prod-br
ip link set veth-dev netns dev
ip link set veth-dev-br master v-net-0
ip link set veth-prod netns prod
ip link set veth-prod-br master v-net-0
ip -n dev addr add 192.168.20.1/24 dev veth-dev
ip -n prod addr add 192.168.20.2/24 dev veth-prod
ip -n dev link set veth-dev up
ip -n prod link set veth-prod up
ip link set veth-dev-br up
ip link set veth-prod-br up
ip -n dev link set lo up
ip -n prod link set lo up

# This needed for to ping from dev to prod and prod to dev.
iptables -A INPUT -i v-net-0 -j ACCEPT
iptables -A FORWARD -i v-net-0 -j ACCEPT