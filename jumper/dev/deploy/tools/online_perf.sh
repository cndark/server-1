#!/bin/bash

FROM=$1
TO=$2

if [ ! $FROM ]; then
    FROM=1
    TO=999
elif [ ! $TO ]; then
    TO=$FROM
fi

for ((i = $FROM; i <= $TO; i++)); do
    svr=game$i

    host=$(../core/deploy.js -h $svr)
    [ $? -ne 0 ] && continue

    ssh game@$host "
        echo -e '============================= \033[32m$svr\033[0m ============================='
        grep 'perfmon' $WORK_DIR/$svr/$svr.log | tail -1
        ps x opcpu,rss,cmd | egrep '(game$i|gate$i)' | grep -v grep | awk '\$2=\$2\"000\"' | numfmt --to=iec --field 2
    "
done

echo -e "\n\n\n\n\n"
