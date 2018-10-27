#!/bin/sh

echo build
go build

echo kill
ps -ef | grep -E "discord-bot" | grep -v grep | awk '{ print "kill  -9", $2 }' | sh

echo start
# nohup.outを出力しない
nohup ./discord-bot > /dev/null 2>&1 &

echo finish
