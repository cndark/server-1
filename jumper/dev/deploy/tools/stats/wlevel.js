#!/usr/bin/env node

require('../../lib/ext');

var path  = require('path');
var dbpool = require('../../lib/dbpool');

// ============================================================================

// var dict = new Map();

// require('./gamedata/data/currency.json').forEach(v => dict.set(v.id, v.name));
// require('./gamedata/data/item.json').forEach(v => dict.set(v.id, v.name));
// require('./gamedata/data/hero.json').forEach(v => dict.set(v.id, v.name));

// ============================================================================

(async () => {
    // open db center
    let str = await shell_exec(`${path.join(__dirname, '../../core/deploy.js')} -c`);
    let c   = JSON.parse(str);

    console.log("uid|关卡id");
    await dbpool.exec(c.common.db_stats, async db => {
        let docs = await db.collection('wlevel').aggregate([
            {$match: {
                svr: "game2",
                // "create_ts":{$lt: new Date('2020-10-18T16:00:00.000Z')},
            }},
            // {$project: {
            //     _id:  0,
            //     wlv: 1,
            // }},
            // {$group:{
            //     _id: `$wlv`,
            //     n: {$sum: 1},
            // }},
        ]).toArray();       

        docs.forEach(row => {

            console.log(row.uid, row.wlv);
        });
    });   
     
})().catch(console.error);
