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

    console.log("订单号 用户id 初心id 金额 产品id cs透传 时间");
    await dbpool.exec(c.common.db_bill, async db => {
        let docs = await db.collection('order').aggregate([
            {$match: {
                "create_ts":{$lt: new Date('2020-10-18T16:00:00.000Z')},
            }},
        ]).toArray();       

        docs.forEach(row => {
            console.log(row.orderid, row.userid, row.chuxin_uid, row.amount, row.prod_id, row.csext ? row.csext : "csext", row.create_ts.toString());
        });
    });

    
     
})().catch(console.error);
