#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

_usage() {
    echo "Usage: $0 cmd [args...]"
    echo
    echo "cmd:"
    echo "    update  svr        update  svr"
    echo "    start   svr        start   svr"
    echo "    stop    svr        stop    svr"
    echo "    restart svr        restart svr"
    echo "    status  svr        show    svr status"
    echo "    reload  svr        reload  svr"
    echo
    exit 1
}

_update() {
    # args
    svr=$1
    [ -z "$svr" ] && _usage

    # query server host
    host=$(core/deploy.js -h $svr)
    [ $? -ne 0 ] && echo "$host" && exit 1

    # check server running state
    ssh game@$host "
        if [ -f $WORK_DIR/$svr/ctl.sh ]; then
            cd $WORK_DIR/$svr
            ./ctl.sh status | grep Running
            exit \$((1 - \$?))
        else
            mkdir -p $WORK_DIR
        fi
    "
    [ $? -ne 0 ] && echo "$svr is still running. shut it down before updating" && exit 1

    # create temp dir
    echo "preparing ..."
    tmpdir=$(mktemp -d ./tmp.$svr.XXXXX)
    svrdir=$tmpdir/$svr
    mkdir $svrdir

    # generate core files
    echo "generate core files ..."
    ./core/deploy.js -c      > $svrdir/config.json
    ./core/deploy.js -s $svr > $svrdir/SERVERS
    ./core/deploy.js -f $svr > $svrdir/FILES

    # copy server files
    echo "copying $svr files ..."
    source $svrdir/FILES

    for f in $FILES; do
        cp -r $BIN/$f $svrdir/
    done
    cp $BIN/ctl.sh $svrdir/

    # pack
    echo "packing ..."
    cd $tmpdir
    tar -czf $svr.tar.gz $svr
    cd ..

    # send to host
    echo "sending to $host ..."
    scp $tmpdir/$svr.tar.gz game@$host:$WORK_DIR

    # update
    echo "updating $svr ..."
    ssh game@$host "
        cd $WORK_DIR

        if [ -d $svr ]; then
            source $svr/FILES
            for f in \$FILES; do
                rm -rf $svr/\$f
            done
        fi

        tar -xf $svr.tar.gz
        rm -rf $svr.tar.gz
    "

    # cleanup
    rm -rf $tmpdir

    echo "Done"
}

_start() {
    # args
    svr=$1
    [ -z "$svr" ] && _usage

    # query server host
    host=$(core/deploy.js -h $svr)
    [ $? -ne 0 ] && echo "$host" && exit 1

    # check server running state
    ssh game@$host "
        if [ -f $WORK_DIR/$svr/ctl.sh ]; then
            cd $WORK_DIR/$svr
            ./ctl.sh start
        fi
    "
}

_stop() {
    # args
    svr=$1
    [ -z "$svr" ] && _usage

    # query server host
    host=$(core/deploy.js -h $svr)
    [ $? -ne 0 ] && echo "$host" && exit 1

    # check server running state
    ssh game@$host "
        if [ -f $WORK_DIR/$svr/ctl.sh ]; then
            cd $WORK_DIR/$svr
            ./ctl.sh stop
        fi
    "
}

_status() {
    # args
    svr=$1
    [ -z "$svr" ] && _usage

    # query server host
    host=$(core/deploy.js -h $svr)
    [ $? -ne 0 ] && echo "$host" && exit 1

    # check server running state
    ssh game@$host "
        if [ -f $WORK_DIR/$svr/ctl.sh ]; then
            cd $WORK_DIR/$svr
            ./ctl.sh status
        fi
    "
}

_reload() {
    # args
    svr=$1
    [ -z "$svr" ] && _usage

    # query server host
    host=$(core/deploy.js -h $svr)
    [ $? -ne 0 ] && echo "$host" && exit 1

    # check server running state
    ssh game@$host "
        if [ -f $WORK_DIR/$svr/ctl.sh ]; then
            cd $WORK_DIR/$svr
            ./ctl.sh reload
        fi
    "
}

# -----------------------------------------------

# set bin dir: production or dev
[ -d ../bin ] && BIN=../bin || BIN=../../../bin

# args
cmd=$1
shift

# go
case $cmd in
    update)
        _update $@
        ;;

    start)
        _start $@
        ;;

    stop)
        _stop $@
        ;;

    restart)
        _stop $@
        _start $@
        ;;

    status)
        _status $@
        ;;

    reload)
        _reload $@
        ;;

    *)
        _usage
        ;;
esac
