#!/bin/bash

source ./env

# =============================================================================
cd $(dirname $0)
# =============================================================================

[ ! $2 ] && echo "$0 <stop|update|status> <svr|exes|webs|all>" && exit 1

ADMIN="./jumper/dev/deploy/admin.sh"

WEBS="auth bill switcher reporter agent admin bat1"
EXES="router1 game1 game2 game3 gate1 gate2 gate3"

case "$2" in
    all)
        echo $WEBS $EXES | xargs -P 25 -n 1 ./dev.sh $1
        ;;

    webs)
        echo $WEBS | xargs -P 25 -n 1 ./dev.sh $1
        ;;

    exes)
        echo $EXES | xargs -P 25 -n 1 ./dev.sh $1
        ;;

    *)
        case "$1" in
            stop)
                $ADMIN stop $2
                ;;

            update)
                $ADMIN stop $2
                $ADMIN update $2
                $ADMIN start $2
                ;;

            status)
                $ADMIN status $2
                ;;
        esac
esac
