#!/bin/bash

cmd=$1
fn=./dict.txt
tardir=../../bin/admin

[ ! $cmd ] && echo "$0 find|replace" && exit 1

# ===================================

_find() {
    declare -A arr

    for e in $(grep --include=*.js --include=*.pug -PRoh '\p{Han}+' $tardir); do
        arr[$e]=1
    done

    printf "%s\n" "${!arr[@]}" > $fn

    echo "chinese words are saved in $fn"
}

_replace() {
    [ -f "$tardir/TRANS.DONE" ] && echo "Already done" && exit 1

    declare -A dict
    while read -r a b; do
        [ ! "$a" ] && continue
        [ ! "$b" ] && echo "$a is NOT translated" && exit 1
        dict[$a]=$b
    done < $fn

    while IFS=':' read -r fn2 a; do
        sed -i "s/$a/${dict[$a]}/g" $fn2
    done <<EOF
$(grep --include=*.js --include=*.pug -PRoH '\p{Han}+' $tardir | awk -F: '{print length($2) " " $0}' | sort -nr | cut -d ' ' -f 2-)
EOF

    touch "$tardir/TRANS.DONE"
    echo "Done."
}

# ===================================

cd $(dirname $0)
fn=$(realpath $fn)

case "$cmd" in
    find)
        #_find
        ;;
    replace)
        _replace
        ;;
esac