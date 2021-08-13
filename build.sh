#!/bin/bash

export GOBIN=$(pwd)/bin

source ./env

# -----------------------------------------------

echo installing router ...
go install ${PROJ_NAME}/src/router

echo installing bat ...
rm -rf bin/bat
cp -R src/bat bin/

echo installing game ...
go install "$@" ${PROJ_NAME}/src/game

echo installing gate ...
go install ${PROJ_NAME}/src/gate

echo installing auth ...
rm -rf bin/auth
cp -R src/auth bin/

echo installing bill ...
rm -rf bin/bill
cp -R src/bill bin/

echo installing switcher ...
rm -rf bin/switcher
cp -R src/switcher bin/

echo installing reporter ...
rm -rf bin/reporter
cp -R src/reporter bin/

echo installing agent ...
rm -rf bin/agent
cp -R src/agent bin/

echo installing admin ...
rm -rf bin/admin
cp -R src/admin bin/

echo installing bot ...
go install ${PROJ_NAME}/src/bot

# -----------------------------------------------

echo preparing node_modules ...
cp -R src/node_modules bin/

echo preparing data ...
cp src/config.json bin/

[ -d bin/gamedata ] && rm -rf bin/gamedata/* || mkdir bin/gamedata
cp -R src/game/app/gamedata/data   bin/gamedata
cp -R src/game/app/gamedata/filter bin/gamedata

echo preparing bot data ...
cp src/bot/botconf/bot.json  bin/
cp src/bot/app/data/chat.txt bin/

# -----------------------------------------------

echo setup server topology

svrfile=bin/SERVERS

printf "\n" > $svrfile
printf "SERVERS=(\n" >> $svrfile

printf "\"%-12s ./auth/www\"\n"                                                    "auth"     >> $svrfile
printf "\"%-12s ./bill/www\"\n"                                                    "bill"     >> $svrfile
printf "\"%-12s ./switcher/www\"\n"                                                "switcher" >> $svrfile
printf "\"%-12s ./reporter/www\"\n"                                                "reporter" >> $svrfile
printf "\"%-12s ./agent/www\"\n"                                                   "agent"    >> $svrfile
printf "\"%-12s ./admin/www\"\n"                                                   "admin"    >> $svrfile
printf "\"%-12s ./router -config config.json -server router1 -log router1.log\"\n" "router1"  >> $svrfile
printf "\"%-12s ./bat/www                    -server bat1                    \"\n" "bat1"     >> $svrfile
printf "\"%-12s ./game   -config config.json -server game1   -log game1.log  \"\n" "game1"    >> $svrfile
printf "\"%-12s ./gate   -config config.json -server gate1   -log gate1.log  \"\n" "gate1"    >> $svrfile

printf ")\n" >> $svrfile

# -----------------------------------------------

echo preparing control script
cp ctl.sh bin

# -----------------------------------------------

echo -e "\033[32mDone\033[0m"
