#!/bin/bash
etcd -name etcd-hc -data-dir /var/lib/etcd http://localhost:2379,http://127.0.0.1:2379 --listen-client-urls http://localhost:2379,http://127.0.0.1:2379 & 
etcdctl  --endpoints http://localhost:2379 member list

etcdctl set  /coreos.com/network/config '{"Network": "10.0.0.0/16", "SubnetLen": 24, "SubnetMin": "10.0.10.0","SubnetMax": "10.0.20.0", "Backend": {"Type": "vxlan"}}'
flannel -etcd-endpoints "http://localhost:4001,http://localhost192.168.1.13:2379"  &
