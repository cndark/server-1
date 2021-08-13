#!/usr/bin/env node

require('../../lib/ext');

var path  = require('path');
var dbpool = require('../../lib/dbpool');

// ============================================================================

// var dict = new Map();

// require('./gamedata/data/currency.json').forEach(v => dict.set(v.id, v.name));
// require('./gamedata/data/item.json').forEach(v => dict.set(v.id, v.name));
// require('./gamedata/data/monster.json').forEach(v => dict.set(v.id, v.color));

// ============================================================================

var hero_ids  = [];

// ============================================================================

(async () => {
    // open db center
    let str = await shell_exec(`${path.join(__dirname, '../../core/deploy.js')} -c`);
    let c   = JSON.parse(str);

    console.log("游戏id|鬼神(id,星级;id,星级;)");
    for (const key in c.common.db_user) {
        if (c.common.db_user.hasOwnProperty(key)) {
            await dbpool.exec(c.common.db_user[key], async db => {
                let docs = await db.collection('user').aggregate([
                    {$match: {
                    }},
                    {$project: {
                        '_id':               1,
                        'base.bag.heroes':   1,
                    }},
                ]).toArray();       
        
                docs.forEach(user => {
                    let b = user.base;
                    let str = "";
                    str += user._id + '|';
 
                    if (b.bag.heroes){
                        Object.keys(b.bag.heroes).forEach(v=>{
                            str += b.bag.heroes[v].id + ',' + b.bag.heroes[v].star + ';'
                        });
                    }

                    console.log(str);
                });
            });
        }
    }
     
})().catch(console.error);
