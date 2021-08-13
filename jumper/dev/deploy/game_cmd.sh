#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

CMD=$1
FROM=$2
TO=$3

if [[ "$CMD" != "start" && "$CMD" != "stop" && "$CMD" != "update" || ! $FROM ]]; then
    echo "$0 <start|stop|update> <all|from-id> [to-id]"
    exit 1
fi
[ ! $TO ] && TO=$FROM

if [ "$FROM" == "all" ]; then
    RANGE="all"
else
    RANGE=$FROM-$TO
fi

read -p "are you sure to $CMD GAMES ($RANGE)? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

# gen game list
if [ "$FROM" == "all" ]; then
    GAMES=$(core/deploy.js -n games)
else
    for ((N=$FROM; N<=$TO; N++)); do
        core/deploy.js -h game$N &>/dev/null
        [ $? -ne 0 ] && continue

        GAMES="$GAMES game$N"
    done
fi

[ ! "$GAMES" ] && echo "no valid game found" && exit 1

# execute cmd
echo "===============> executing $CMD on games ..."
echo $GAMES | xargs -P 25 -n 1 ./admin.sh $CMD

echo "===============> done."
