#!/bin/bash

[ ! $2 ] && echo "$0 svr batch" && exit 1

./genconf.js $1 $2 > bot-$1-$2.json
./bot -config bot-$1-$2.json &> $1-$2.log &
