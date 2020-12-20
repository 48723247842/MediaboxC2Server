#!/bin/bash
rm go.mod
rm go.sum

sudo rm -rf ~/go/pkg/mod/github.com/0187773933/ || echo ""
sudo rm -rf  /home/morphs/go/src/$1 || echo ""
sudo rm -rf  /home/morphs/go/pkg/mod/$1 || echo ""
go clean --modcache

go mod init c2server

# We have to force golang to not care about cache of repos we are currently updating
RMUHash=$(curl -s 'https://api.github.com/repos/0187773933/RedisManagerUtils/git/refs/heads/master' | jq -r '.object.sha')
go get "github.com/0187773933/RedisManagerUtils/@$RMUHash"

go run main.go