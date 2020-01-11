#!/usr/bin/env bash

PID=$(cat ./pid_num)
echo $PID
running=$(ps aux | grep $PID)
if [ "$running" != "" ]
then
    kill $PID
    echo "stop "
else
    echo "not running "
fi
