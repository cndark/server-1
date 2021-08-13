#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

# args:
#   none    output the next wait id
#   all     output all wait ids

Ids=($(awk '/^\s*\/\/\s*wait\s*{/   {match($0, /id\s*:\s*([[:digit:]]+)/, a); print a[1]}' 3-servers.js))
[ "$1" == "all" ] && echo "${Ids[*]}" || echo ${Ids[0]}
