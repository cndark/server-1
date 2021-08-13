#!/usr/bin/env node

require('../lib/ext');

var path = require('path');

var dbpool = require('../lib/dbpool');

// ============================================================================

(async () => {
    // scan area dirs
    let p = path.join(__dirname, '../../../');
    let dirs = await shell_exec(`bash -c '
        for f in $(find "${p}" -maxdepth 1 -type d -iregex ".*[^./]"); do
            if [ -d $f/deploy ]; then
                echo $f
            fi
        done
    '`);
    dirs = dirs.split('\n').filter(v => v.trim().length > 0);

    // open db share
    let str = await shell_exec(`${path.join(__dirname, './deploy.js')} -c`);
    let c   = JSON.parse(str);

    await dbpool.exec(c.common.db_share, async db => {
        // remove all
        await db.collection('areas').deleteMany({});

        // insert
        await Promise.all(dirs.map(async dir => {
            let str  = await shell_exec(`${dir}/deploy/core/deploy.js -c`);
            let c    = JSON.parse(str);
            let area = c.common.area;

            let doc = {
                _id:    area.id,
                name:   area.name,
                dir:    dir,
                config: c,
            };

            await db.collection('areas').insertOne(doc);
        }));
    });
})().catch(e => {
    console.error(typeof e == 'string' ? e : e.message);
    process.exit(1);
});
