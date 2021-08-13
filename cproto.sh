#!/bin/bash

source ./env

# retrive code rev
echo "retrieving code version ..."
branch_code=$(git branch | grep \* | awk '{print $2}')

# read branch config
branch_proto_var=${branch_code}_branch_proto
branch_proto=${!branch_proto_var}
[ ! $branch_proto ] && branch_proto=$branch_code

echo "    proto branch: $branch_proto"

# dirs
case $LOGNAME in
    evil)
        SERVER_DIR=~/${PROJ_NAME}-hserver
        SHARED_DIR=~/shared/${PROJ_NAME}
        ;;

    haa|xuezhang)
        SERVER_DIR=~/${PROJ_NAME}-hserver
        SHARED_DIR=~/${PROJ_NAME}-hshared
        ;;

    *)
        echo "who the hell are you?"
esac

# update
svn update ${SHARED_DIR}/server/${branch_proto}

rm -rf ${SHARED_DIR}/server/${branch_proto}/proto/server
cp -r ${SERVER_DIR}/src/proto ${SHARED_DIR}/server/${branch_proto}/proto/server

cd ${SHARED_DIR}/server/${branch_proto}/proto
lua ../../gen_TS.lua
cd - > /dev/null

while read -r a b; do
    if [ "$a" == '!' ]; then
        svn rm $b
    elif [ "$a" == '?' ]; then
        svn add $b
    fi
done <<EOF
$(svn status ${SHARED_DIR}/server/${branch_proto}/proto)
EOF

svn commit -m 'proto: updated' ${SHARED_DIR}/server/${branch_proto}/proto
