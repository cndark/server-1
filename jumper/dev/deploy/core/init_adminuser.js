#!/usr/bin/env node

require('../lib/ext');

var path = require('path');
var md5  = require('md5');

var dbpool = require('../lib/dbpool');

// ============================================================================

let args = process.argv.slice(2);
let pwd = args.length > 0 ? args[0] : 'default888';

// ============================================================================

(async () => {
    let str = await shell_exec(`${path.join(__dirname, './deploy.js')} -c`);
    let c   = JSON.parse(str);

    await dbpool.exec(c.common.db_share, async db => {
        await db.collection('adminuser').replaceOne(
            {_id: 'admin'},
            {pwd: md5(pwd + c.admin.pwdfill), rank: 0, priv: 'all', memo: 'super user'},
            {upsert: true},
        );
    });
})().catch(e => {
    console.error(typeof e == 'string' ? e : e.message);
    process.exit(1);
});
