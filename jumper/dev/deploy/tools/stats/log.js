#!/usr/bin/env node

require('../../lib/ext');

var path = require('path');
var dbpool = require('../../lib/dbpool');

// ============================================================================

// var dict = new Map();

// require('./gamedata/data/currency.json').forEach(v => dict.set(v.id, v.name));
// require('./gamedata/data/item.json').forEach(v => dict.set(v.id, v.name));
// require('./gamedata/data/monster.json').forEach(v => dict.set(v.id, v.color));

// ============================================================================

// ============================================================================

(async () => {
    // open db center
    let str = await shell_exec(`${path.join(__dirname, '../../core/deploy.js')} -c`);
    let c = JSON.parse(str);

    // console.log("游戏id");
    await dbpool.exec(c.common.db_log, async db => {
        let docs = await db.collection('log').aggregate([
            {
                $match: {
                    'role_id': { $in: ["u2-101-1000521", "u1-101-1000230", "u2-101-1000940"] },
                    'func': { $in: ["moneyChange", "bagChange"] },

                    'action_type': { $in: [100, 101, 102, 103, 104, 800, 801, 802, 803, 804, 805, 806, 807, 808] },
                }
            },
            {
                $project: {
                    '_id': 0,
                    'finish_time': 0,
                    'channel_id': 0,
                    'role_create_time': 0,
                    'platform_id': 0,
                    'server_id': 0,
                    'account': 0,
                    'ip': 0,
                    'timestamp': 0,
                    'device_id': 0,
                    'cts': 0,
                }
            },
        ]).toArray();

        docs.forEach(doc => {
            if (doc.func == "moneyChange") {
                console.log(doc.role_id, doc.ts, doc.status, doc.action_type, doc.currency_type, doc.currency_num, doc.after_num);
            } else {
                console.log(doc.role_id, doc.ts, doc.status, doc.action_type, doc.goods_id, doc.goods_num, doc.goods_residue);
            }
        });
    });


})().catch(console.error);
