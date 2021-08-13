#!/bin/bash

# =============================================================================
cd $(dirname $0)
# =============================================================================

echo "edit first!"
exit 1

read -p "!!!!!!!! Are you sure?[Yes!Iam./No]" x
[ "$x" != "Yes!Iam." ] && exit 1

echo "make sure you know what you're doing!!!"
exit 1

for d in db_xxx; do
    cnnstr=$(core/deploy.js -d $d)
    mongo $cnnstr --quiet --eval "
        var cnn = db.getMongo();
        cnn.getDBNames().forEach(name => {
            if (name.startsWith('$PROJ_NAME')) {
                cnn.getDB(name).getCollectionNames().forEach(coll => {
                    var r = cnn.getDB(name).getCollection(coll).drop();
                    printjson(r);
                });
            }
        });
    "
done

