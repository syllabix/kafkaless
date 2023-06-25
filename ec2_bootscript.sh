#!/bin/bash
yum update -y
yum install -y git
cd /tmp
wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go*.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
go install github.com/ServiceWeaver/weaver/cmd/weaver@latest