#!/usr/bin/env bash

echo "starting app by $(whoami)"

nohup ./main > main.log &
sleep 1 # just make console log order more normal
app_pid=$!
echo $app_pid > pid_num

echo 'start success'
