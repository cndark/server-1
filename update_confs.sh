#!/bin/bash

read -p "continue to update confs? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

# -------------------------------------

set -e

source ./env

# retrive code rev
echo "retrieving code version ..."
branch_code=$(git branch | grep \* | awk '{print $2}')

# read branch config
branch_data_var=${branch_code}_branch_data
branch_bat_var=${branch_code}_branch_bat

branch_data=${!branch_data_var}
branch_bat=${!branch_bat_var}

[ ! $branch_data ] && branch_data=$branch_code
[ ! $branch_bat  ] && branch_bat=$branch_code

echo "    data branch: $branch_data"
echo "    bat  branch: $branch_bat"

# update gamedata
echo "updating gamedata ..."
rm src/game/app/gamedata/data   -rf
rm src/game/app/gamedata/filter -rf

svn export --force "${SVN_GAMEDATA_URL}_${branch_data}/go"     src/game/app/gamedata > /dev/null
svn export --force "${SVN_GAMEDATA_URL}_${branch_data}/filter" src/game/app/gamedata/filter > /dev/null

REV_GAMEDATA=$(svn export --force "${SVN_GAMEDATA_URL}_${branch_data}/json" src/game/app/gamedata/data|tail -1)
[[ "$REV_GAMEDATA" != "Exported revision"* ]] && exit 1
echo "    gamedata: $REV_GAMEDATA"

# update calcbattle
echo "updating bat code ..."
rm src/bat/calcbattle -rf
REV_CALCBATTLE=$(svn export --force "${SVN_CALCBATTLE_URL}/${branch_bat}/bat/calcbattle" src/bat/calcbattle|tail -1)
[[ "$REV_CALCBATTLE" != "Exported revision"* ]] && exit 1
echo "    calcbattle: $REV_CALCBATTLE"

# write version file
echo "preparing version file ..."
cat > .ver/VER_${branch_code} << EOF
gamedata:
    $branch_data
    $REV_GAMEDATA
calcbattle:
    $branch_bat
    $REV_CALCBATTLE
EOF

# done
echo -e "\033[32mDone.\033[0m"
