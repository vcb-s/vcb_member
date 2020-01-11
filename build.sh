#!/usr/bin/env bash

echo "begin..."

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build main.go
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build main.go

echo "upload..."
scp build -P 30000 -2 inori@bupt.moe:/www/vcb-s/vcbs_member/be/

echo "done"
