#!/usr/bin/env bash

PID=$(cat ./pid_num)
echo "finding process $PID status"
if ps -p $PID > /dev/null
then
    echo "$PID is running."
    kill "$PID"
    echo "$PID stoped."
else
    echo "not running"
fi
