#!/usr/bin/env node

require('../../lib/ext');

var path = require('path');
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
    let c = JSON.parse(str);

    console.log("账户id|游戏id|创号时间|累计登陆天数|在线时长秒|上次离线时间|ip|等级|充值金额|vip|关卡停留|好友个数|剩余钻石");
    for (const key in c.common.db_user) {
        if (c.common.db_user.hasOwnProperty(key)) {
            await dbpool.exec(c.common.db_user[key], async db => {
                let docs = await db.collection('user').aggregate([
                    {
                        $match: {
                            // 'base.create_ts':{$lt: new Date('2020-10-16T16:00:00.000Z')},
                        }
                    },
                ]).toArray();

                docs.forEach(user => {
                    let b = user.base;
                    let str = "";
                    str += b.authid + '|';
                    str += user._id + '|';
                    str += b.create_ts.toString() + '|';
                    str += b.login_sumdays + '|';
                    str += b.online_dur + '|';
                    str += b.off_ts.toString() + '|';
                    str += b.login_ip + '|';
                    str += b.lv + '|';
                    str += b.bill.totalbaseccy + '|';
                    str += b.vip.lv + '|';
                    str += '|';
                    str += b.wlevel.lvnum + '|';

                    if (b.friend.frds) {
                        str += Object.keys(b.friend.frds).length;
                    }
                    str += '|';

                    if (b.bag.ccy) {
                        str += b.bag.ccy['10002'] ? b.bag.ccy['10002'] : 0;
                    }

                    console.log(str);
                });
            });
        }
    }

})().catch(console.error);
