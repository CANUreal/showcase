#!/bin/bash
set -e

echo "configuring firewalld for k3s"
if systemctl is-active --quiet firewalld; then
    sudo firewall-cmd --zone=trusted --add-interface=cni0 --permanent
    sudo firewall-cmd --zone=public --add-port=8472/udp --permanent
    sudo firewall-cmd --reload
    echo "firewalld rules applied"
else
    echo "firewalld is not active, skipping"
fi
