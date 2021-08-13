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

    console.log("类型|步数|人数");
    await dbpool.exec(c.common.db_stats, async db => {
        let A = [
            ["1","新手"],
            ["2","妖兽"],
            ["3","排行榜"],
            ["4","首次失败"],
            ["5","任务"],
            ["6","猎魂"],
            ["7","限时冲榜"],
            ["8","结界修炼"],
            ["9","家族"],
            ["10","祈愿"],
            ["11","鬼神试炼"],
            ["12","灵域争夺"],
            ["13","极乐之塔"],
            ["14","自动战斗"],
            ["15","妖界历练"],
            ["16","鬼神委派"],
            ["17","统治地狱"],
            ["18","器灵"],
            ["19","宴会"],
            ["20","合魂"],
        ];

        for (let index = 0; index < A.length; index++) {
            const element = A[index];
            
            let t = `$tut.` + element[0];
            let docs = await db.collection('tutorial').aggregate([
                {$match: {
                    svr: "game2",
                    // "create_ts":{$lt: new Date('2020-10-18T16:00:00.000Z')},
                }},
                {$project: {
                    _id:  0,
                    tp: t,
                }},
                {$group:{
                    _id: `$tp`,
                    n: {$sum: 1},
                }},
            ]).toArray();       

            docs.forEach(row => {
                console.log(element[1], row._id ? row._id : 0, row.n)
            });
        }
    });

    
     
})().catch(console.error);
