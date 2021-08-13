#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

read -p "are you sure to init deployment? [Yes/No]" x
[ "$x" != "Yes" ] && exit 1

../core/update_areas.js
../core/init_adminuser.js
