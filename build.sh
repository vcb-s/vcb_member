#!/usr/bin/env bash

#linux  64位
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o index index.go

# scp index inori@vcb-s.com:/www/vcbs_member/be/index.new
