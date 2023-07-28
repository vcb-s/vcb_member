#!/usr/bin/env bash

echo "begin..."

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build main.go
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build main.go

echo "build success"

echo "uploading..."
# windows bash 环境下的scp实在多少有点问题…会只剩下100k不到；linux/mac下边未遇见类似问题
scp build vcb-s:/var/www/vcb-s/vcbs_member/be/build
echo "upload success"

echo "restarting..."
echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>"
ssh -t vcb-s "cd /var/www/vcb-s/vcbs_member/be/ && . restart.sh"
echo "<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
echo "restart success"

echo "deploy success"
