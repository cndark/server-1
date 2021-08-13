#!/bin/bash

[ ! $1 ] && echo "$0 backup|restore" && exit 1

read -p "continue? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

REPO=$(realpath ../../db-gift)
mkdir -p $REPO

DB=($(node -e '
var dbs = require("../1-dbs");
for (let db of dbs) {
    if (db.name == "db_c") {
        console.log(db.ip, db.port, db.user, db.pwd);
        break;
    }
}
'))

IP=${DB[0]}
PORT=${DB[1]}
USER=${DB[2]}
PWD=${DB[3]}

echo "===========> $IP, $PORT"

case $1 in
    backup)
        mongodump -h $IP --port $PORT -u $USER -p "$PWD" \
            --authenticationDatabase=admin \
            -d ${PROJ_NAME}s -c giftinfo --gzip --archive=$REPO/giftinfo.db
        ;;

    restore)
        mongorestore -h $IP --port $PORT -u $USER -p "$PWD" \
            --authenticationDatabase=admin \
            --drop --gzip --archive=$REPO/giftinfo.db
        ;;

    *)
        echo "invalid command!"
esac
