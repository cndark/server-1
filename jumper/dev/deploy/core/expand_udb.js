#!/usr/bin/env node

require('../lib/ext');

var path = require('path');

var dbpool = require('../lib/dbpool');

// ============================================================================

(async () => {
    // open db center
    let str = await shell_exec(`${path.join(__dirname, './deploy.js')} -c`);
    let c   = JSON.parse(str);

    await dbpool.exec(c.common.db_center, async db => {
        for (let k in c.common.db_user) {
            try {
                await db.collection('userload').insertOne({_id: k, n: 0});
            } catch(e) {
                if (!e.message.match(/^E11000/)) throw e;
            }
        }
    });
})().catch(e => {
    console.error(typeof e == 'string' ? e : e.message);
    process.exit(1);
});
