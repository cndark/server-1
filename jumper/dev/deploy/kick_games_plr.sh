#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

FROM=$1
TO=$2

if [[ ! $FROM ]]; then
    echo "$0 <from-id> [to-id]"
    exit 1
fi
[ ! $TO ] && TO=$FROM

RANGE=$FROM-$TO

read -p "are you sure to kick player in game($RANGE)? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

# gen game list
for ((N=$FROM; N<=$TO; N++)); do
    core/deploy.js -h game$N &>/dev/null
    [ $? -ne 0 ] && continue

    GAMEIDS="$GAMEIDS $N"
done

[ ! "$GAMEIDS" ] && echo "no valid game found" && exit 1

# execute
for gate in $(core/deploy.js -n gates); do
    echo "handle $gate ..."

    gate_dir=$WORK_DIR/$gate
    ssh game@$(core/deploy.js -h $gate) "
        if [ -f $gate_dir/$gate.pid ]; then
            echo $GAMEIDS > $gate_dir/KICK_GAMEIDS

            pid=$(< $gate_dir/$gate.pid)
            kill -s USR1 \$pid
        fi
    "
done

echo "===============> done."
