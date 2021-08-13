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

var ccy_ids  = [];
var item_ids = ["20181","20003","20004","20020","20021","20022","20087","20166","20008"];

// ============================================================================

(async () => {
    // open db center
    let str = await shell_exec(`${path.join(__dirname, '../../core/deploy.js')} -c`);
    let c   = JSON.parse(str);

    console.log("游戏id|物品(物品id,n;物品id,n;)");
    for (const key in c.common.db_user) {
        if (c.common.db_user.hasOwnProperty(key)) {
            await dbpool.exec(c.common.db_user[key], async db => {
                let docs = await db.collection('user').aggregate([
                    {$match: {
                    }},
                    {$project: {
                        '_id':            1,
                        'base.bag.ccy':   1,
                        'base.bag.items': 1,
                    }},
                ]).toArray();       
        
                docs.forEach(user => {
                    let b = user.base;
                    let str = "";
                    str += user._id + '|';
 
                    if (b.bag.ccy){
                        ccy_ids.forEach(id=>{
                            if (b.bag.ccy[id]){
                                str += id + ',' + b.bag.ccy[id] + ';';
                            }
                        });
                    }

                    if (b.bag.items){
                        item_ids.forEach(id=>{
                            if (b.bag.items[id]){
                                str += id + ',' + b.bag.items[id] + ';';
                            }
                        });
                    }

                    console.log(str);
                });
            });
        }
    }
     
})().catch(console.error);
