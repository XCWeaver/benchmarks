#!/bin/bash

# temporarily disabled because it is getting stuck
#sudo apt update -y && sudo apt upgrade -y
sudo apt install -y docker.io docker-compose dnsutils curl wget rsync git tmux python3-pip

# install Go 1.21.5
sudo wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz && sudo rm -rf go1.21.5.linux-amd64.tar.gz
# by default the Go binary is placed in /usr/local/go/bin and
# binaries of Go modules (e.g. Weaver) are placed in $HOME/go/bin
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# install Weaver 0.24.5
go install github.com/ServiceWeaver/weaver/cmd/weaver@vv0.24.5
export PATH="$PATH:$HOME/go/bin"
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
source ~/.bashrc

#install kubectl
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.31/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
sudo chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.31/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo chmod 644 /etc/apt/sources.list.d/kubernetes.list
sudo apt-get install -y kubectl

# install kube deployer 0.24.7
go install github.com/ServiceWeaver/weaver-kube/cmd/weaver-kube@v0.24.7