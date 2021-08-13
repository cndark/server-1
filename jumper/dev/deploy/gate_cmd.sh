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

read -p "are you sure to $CMD GATES ($RANGE)? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

# gen gate list
if [ "$FROM" == "all" ]; then
    GATES=$(core/deploy.js -n gates)
else
    for ((N=$FROM; N<=$TO; N++)); do
        core/deploy.js -h gate$N &>/dev/null
        [ $? -ne 0 ] && continue

        GATES="$GATES gate$N"
    done
fi

[ ! "$GATES" ] && echo "no valid gate found" && exit 1

# execute cmd
echo "===============> executing $CMD on gates ..."
echo $GATES | xargs -P 25 -n 1 ./admin.sh $CMD

echo "===============> done."
