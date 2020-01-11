#!/usr/bin/env bash

nohup ./main > main.log   &
LOVE_PID=$!
echo $LOVE_PID>pid_num