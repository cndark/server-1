#!/bin/bash

source ./env

IP=$1
AREA_DIR=$2
shift; shift
DEPLOY_ARGS=$@

[ ! $IP ] && echo "$0 ip area_dir [deploy-args...]" && exit 1

read -p "continue to release? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

# -------------------------------------

set -e

# target
host=game@$IP

# build
go clean -cache
mv bin bin_old
./build.sh

# retrive code rev
echo "retrieving code version ..."
branch_code=$(git branch | grep \* | awk '{print $2}')
COMMIT=$(git log | head -1 | awk '{print $2}')

# copy version file
cp .ver/VER_${branch_code} bin/
cat >> bin/VER_${branch_code} << EOF
code:
    $branch_code $COMMIT
EOF

# pack
echo "packing ..."
tarfile=${PROJ_NAME}-${branch_code}.$(date +%Y%m%d.%H%M).tar.gz
tar -czf $tarfile bin

# upload
echo "uploading ..."
scp $tarfile $host:$WORK_DIR/$AREA_DIR

# clean up
echo "cleaning up ..."
rm $tarfile
rm -rf bin
mv bin_old bin

# exec deploy.sh if there's any in remote
ssh $host "
    if [ -x $WORK_DIR/$AREA_DIR/dev_auto_update.sh ]; then
        cd $WORK_DIR/$AREA_DIR
        ./dev_auto_update.sh $tarfile $DEPLOY_ARGS
    fi
"

# done
echo -e "\033[32mDone.\033[0m"
