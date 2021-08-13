#!/bin/bash

set -e

# check
[ ! $1 ] && echo "$0 ids..." && exit 1

# make args
for i in "$@"; do
    games="$games game$i"
done

# gen config.json
echo "=============> generating config.json ..."
../../core/deploy.js -c > config.json

# check
./merge.js --check $games

# stop servers that are to be merged
echo "=============> stopping servers for merging ..."
echo $games | xargs -P 25 -n 1 ../../admin.sh stop

# merge
echo "=============> merging ..."
./merge.js $games

# comment merged entries in 3-servers.js
echo "=============> updating 3-servers.js ..."
for i in "${@:2}"; do
    sed -i "/^\s*games\s*=\s*$/, /^\s*},\s*$/ s/^\s*{\s*id\s*=\s*$i\s*,/---------- M ---------- &/g" ../../3-servers.js
done

# reload
../../reload.sh <<< Yes

# start
echo "=============> starting merged server ..."
../../admin.sh start game$1

# done
echo "=============> Done."
