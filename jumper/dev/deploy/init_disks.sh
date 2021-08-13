#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

names=$@

[ ! $names ] && echo "$0 svr-name..." && exit 1

read -p "continue? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

for x in $names; do
   echo ================== $x
   ssh -t game@$(core/deploy.js -h $x) "sudo /root/mount-vdb.sh"
done
