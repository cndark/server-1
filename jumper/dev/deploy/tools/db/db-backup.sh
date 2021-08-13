#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

read -p "continue to backup? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

REPO=$(realpath ../../../db-backup)
mkdir -p $REPO

./db-cmd.js -d | xargs -P $(./db-cmd.js --nhost) -L 1 bash -c "
    v=\$*
    ip=\${v%% mongo*}
    cmd=\${v#* }
    fn=\${v##*=}

    echo -e \"============================>>> START \$fn ==============================>>>\n\"

    ssh game@\$ip \"
        cd $WORK_DIR
        \$cmd
    \"

    scp game@\$ip:$WORK_DIR/\$fn $REPO

    ssh game@\$ip \"
        cd $WORK_DIR
        rm -rf \$fn
    \"

    echo -e \"<<<============================ END \$fn <<<==============================\n\"
" _

# done
echo -e "\033[32mDone\033[0m"
