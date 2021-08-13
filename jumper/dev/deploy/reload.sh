#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

read -p "are you sure to reload? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

# generate config.json
tmpdir=$(mktemp -d ./tmp.config.XXXXX)
core/deploy.js -c > $tmpdir/config.json

# update config.json
echo "===============> updating config.json ..."
core/deploy.js -n | xargs -P 25 -n 1 bash -c "
    host=\$(core/deploy.js -h \$1)
    scp $tmpdir/config.json game@\$host:$WORK_DIR/\$1
" _

# update area info
echo "===============> updating area info ..."
core/update_areas.js

# reload
echo "===============> reloading ..."
core/deploy.js -n | xargs -P 25 -n 1 ./admin.sh reload

# wait some time
echo "===============> wait ..."
sleep 3

# expand user-db
core/expand_udb.js

# cleanup
rm -rf $tmpdir
