#!/usr/bin/env bash

echo "begin..."

go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.io,direct

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build main.go
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build main.go

# 不使用scp，上传不知道为何很慢很慢，sftp协议有3M/s，这里只有500kbps的平均值
# echo "upload..."
# scp build -P 30000 -2 inori@bupt.moe:/www/vcb-s/vcbs_member/be/

echo "done"
