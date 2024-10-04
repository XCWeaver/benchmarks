#!/bin/bash

# temporarily disabled because it is getting stuck
#sudo apt update -y && sudo apt upgrade -y
sudo apt install -y wget git tmux python3-pip rsync

# install Go 1.21.5
sudo wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz && sudo rm -rf go1.21.5.linux-amd64.tar.gz
# by default the Go binary is placed in /usr/local/go/bin and
# binaries of Go modules (e.g. Weaver) are placed in $HOME/go/bin
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# install XCWeaver 0.22.0
go install github.com/XCWeaver/xcweaver/cmd/xcweaver@v0.22.0
export PATH="$PATH:$HOME/go/bin"
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
source ~/.bashrc

# install Lua Dependencies
sudo apt-get install -y libssl-dev luarocks
sudo luarocks install lua-json
sudo luarocks install luasocket
sudo luarocks install penlight