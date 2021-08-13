#!/bin/bash

read -p "Are you sure to clear dbs? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

source ./env

# get host info
GAME_HOST=game@$(jumper/dev/deploy/core/deploy.js -h game1)
DB_HOST=$(jumper/dev/deploy/core/deploy.js -d db_c)

# kill all games
ssh $GAME_HOST "
    ps aux | grep '\-server game' | grep -v grep > /dev/null
    if [ \$? -eq 0 ]; then
        ps aux | grep '\-server game' | grep -v grep | awk '{print \$2}' | xargs kill -9
    fi
"

# wait
sleep 0.5

# clear
mongo $DB_HOST --quiet --eval "
    var cnn = db.getMongo();
    cnn.getDBNames().forEach(name => {
        if (name.startsWith('${PROJ_NAME}')) {
            cnn.getDB(name).getCollectionNames().forEach(coll => {
                cnn.getDB(name).getCollection(coll).drop();
            });
            print(name + ' cleared');
        }
    });
"

ssh $GAME_HOST "
    # insert dev svrlist
    $WORK_DIR/dev/deploy/core/add_svrlist.js 1 '1服'
    $WORK_DIR/dev/deploy/core/add_svrlist.js 2 '2服'
    $WORK_DIR/dev/deploy/core/add_svrlist.js 3 '3服'

    # update area info
    $WORK_DIR/dev/deploy/core/update_areas.js

    # add adminuser
    $WORK_DIR/dev/deploy/core/init_adminuser.js 1
"
