#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

# 说明：
#   3-servers.js 中注释写法有严格规定 ->
#   game 前 使用 // wait 来注释, 表示后备新服等待中

# args
[ ! $1 ] && echo "$0 'next'|gameid" && exit 1

if [ "$1" != "next" ]; then
    REQ_ID=$1
fi

# scan for new waited game
ID=$(./next_wait_id.sh)
[ ! $ID ] && echo 'no waited game found' && exit 1

# check id consistency if requested
[[ $REQ_ID && "$REQ_ID" != "$ID" ]] && echo "requested id is $REQ_ID, but arranged id should be $ID" && exit 1

# confirm
read -p "are you sure to open NEW game$ID? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

NAME=S$ID

set -e

# 打开 3-servers.js 中的注释
echo "updating server config ..."
sed -i "/^\s*\/\/\s*wait\s*{\s*id\s*:\s*$ID\s*,/   s#//\s*wait\s*##" 3-servers.js

# 更新 game 程序
echo "updating game$ID files ..."
./admin.sh update game$ID

# 开服
echo "starting game$ID ..."
./admin.sh start game$ID
sleep 10

# reload
./reload.sh <<< Yes

# 插入服务器列表数据
echo -e "\033[32mupdating svrlist ...\033[0m"
./core/add_svrlist.js $ID $NAME -s 3 -f new
