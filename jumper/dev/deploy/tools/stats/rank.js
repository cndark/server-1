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

    console.log("本服榜榜id|uid|分数");
    await dbpool.exec(c.common.db_cross, async db => {
        let docs = await db.collection('ranklocal').aggregate([
            {$match: {
                _id: 100104,
                svrid: 1,
            }},
            {$unwind: "$d"},
            {$sort: {
                "d.score": -1,
            }}      
        ]).toArray();       

        docs.forEach(row => {
            console.log(`${row.rankid}|${row.d.info.id}|${row.d.score}`);
        });
    });   
     
})().catch(console.error);
