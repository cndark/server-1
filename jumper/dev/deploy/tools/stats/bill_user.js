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

    console.log("账户id|ip|游戏id|创号时间|上次离线时间|创号天数|在线时长秒|等级|关卡停留|充值金额|vip|第几天登录过|天数_充值次数_金额");
    for (const key in c.common.db_user) {
        if (c.common.db_user.hasOwnProperty(key)) {
            await dbpool.exec(c.common.db_user[key], async db => {
                let docs = await db.collection('user').aggregate([
                    {
                        $match: {
                            // '_id': { $in: ["u1-101-1000230"] },
                            'base.create_ts': { $gt: new Date('2021/7/19 00:00:00'), $lt: new Date('2021/7/20 00:00:00') },
                            'base.bill.totalbaseccy': { $gt: 0 },
                            'base.svr': 'game7',
                        }
                    },
                ]).toArray();

                for (let index = 0; index < docs.length; index++) {
                    const user = docs[index];
                    let b = user.base;
                    let str = "";
                    str += b.authid + '|';
                    str += b.login_ip + '|';
                    str += user._id + '|';
                    str += b.create_ts.toString() + '|';
                    str += b.off_ts.toString() + '|';
                    str += (b.off_ts.startOfDay().unix() - b.create_ts.startOfDay().unix()) / 86400 + '|';
                    str += b.online_dur + '|';
                    str += b.lv + '|';
                    str += b.wlevel.lvnum + '|';
                    str += b.bill.totalbaseccy / 100 + '|';
                    str += b.vip.lv + '|';

                    let loginDays = [];
                    await dbpool.exec(c.common.db_stats, async db => {
                        let docs2 = await db.collection('login').aggregate([
                            {
                                $match: {
                                    'uid': user._id,
                                }
                            },
                            {
                                $project: {
                                    'day': 1,
                                }
                            },
                        ]).toArray();

                        docs2.forEach(v => {
                            loginDays.push(v.day);
                        });
                    });
                    str += loginDays.join(',') + '|';

                    let billDays = [];
                    await dbpool.exec(c.common.db_stats, async db => {
                        let docs2 = await db.collection('bill').aggregate([
                            {
                                $match: {
                                    'uid': user._id,
                                }
                            },
                            {
                                $project: {
                                    'day': 1,
                                    'amt': 1,
                                    'n': 1,
                                }
                            },
                        ]).toArray();

                        docs2.forEach(v => {
                            if (v.amt && v.amt > 0) {
                                billDays.push(v.day + '_' + v.n + '_' + v.amt / 100);
                            }
                        });
                    });
                    str += billDays.join(',') + '|';

                    console.log(str);
                }

            });
        }
    }

})().catch(console.error);