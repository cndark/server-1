#!/bin/bash

read -p "continue to collect client-callstack? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

ip_reporter=$(../core/deploy.js -h reporter)
ip_switcher=$(../core/deploy.js -h switcher)

ssh game@$ip_reporter "
    cd $WORK_DIR
    tar -C reporter -czf ccs.tar.gz ccs
"

scp -3 game@$ip_reporter:$WORK_DIR/ccs.tar.gz game@$ip_switcher:$WORK_DIR/update

ssh game@$ip_reporter "
    cd $WORK_DIR
    rm ccs.tar.gz
"
