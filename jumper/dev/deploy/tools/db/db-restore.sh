#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

read -p "continue to restore? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

echo "edit first!"
exit 1

read -p "proceed? [Yes!I am aware of what i am doing.!./No]" x
[ "$x" != "Yes!I am aware of what i am doing.!." ] && exit 1

REPO=$(realpath ../../../db-backup)

[ ! -d $REPO ] && echo "dir $REPO NOT found" && exit 1

./db-cmd.js -r | xargs -P $(./db-cmd.js --nhost) -L 1 bash -c "
    v=\$*
    ip=\${v%% mongo*}
    cmd=\${v#* }
    fn=\${v##*=}

    echo -e \"============================>>> START \$fn ==============================>>>\n\"

    scp $REPO/\$fn game@\$ip:$WORK_DIR

    ssh game@\$ip \"
        cd $WORK_DIR
        #\$cmd
        rm -rf \$fn
    \"

    echo -e \"<<<============================ END \$fn <<<==============================\n\"
" _

# done
echo -e "\033[32mDone\033[0m"
