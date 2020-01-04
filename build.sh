#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build main.go
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build main.go

# scp build inori@vcb-s.com:/www/vcbs_member/be/main.new
