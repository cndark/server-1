
var express = require('express');
var router  = express.Router();

var dbpool  = require('../../../lib/dbpool');
var gtab    = require('../../../models/gtab');
var stats   = require('../stats');

// ============================================================================

router.get('/', _A_(async (req, res) => {
    var q = stats.make_q(req.query);

    // render func
    var render = (tab) => {
        tab = tab || {};
        tab.name = '付费概览';
        stats.render(req, res, tab);
    }

    // ----------------------

    var tz = new Date().getTimezoneOffset() * 60 * 1000;

    // ----------------------

    try {
        let db = dbpool.get('stats');

        let fs = {
            login: async () => {
                let pipeline = [
                    {$match: {
                        lts: {$gte: q.d_start, $lt: q.d_end},
                    }},
                    {$project: {
                        lts: {$subtract: ["$lts", tz]},
                        uid: 1,
                    }},
                    {$group: {
                        _id:  {$dateToString: {format: "%Y-%m-%d", date: "$lts"}},
                        uids: {$push: "$uid"},
                    }},
                    {$sort: {_id: -1}},
                ];

                stats.filter(pipeline, q);

                return await db.collection('login').aggregate(pipeline).toArray();
            },

            bill: async () => {
                let pipeline = [
                    {$match: {
                        bts: {$gte: q.d_start, $lt: q.d_end},
                        amt: {$gt: 0},
                    }},
                    {$project: {
                        bts: {$subtract: ["$bts", tz]},
                        uid: 1,
                        amt: 1,
                    }},
                    {$group: {
                        _id:  {$dateToString: {format: "%Y-%m-%d", date: "$bts"}},
                        uids: {$push: "$uid"},
                        amt:  {$sum: "$amt"},
                    }},
                    {$sort: {_id: -1}},
                ];

                stats.filter(pipeline, q);

                return await db.collection('bill').aggregate(pipeline).toArray();
            }
        };

        let result = await Promise.all(['login', 'bill'].map(k => fs[k]()));
        result.login = result[0];
        result.bill  = result[1];

        // calc user-count for sum
        var sum_ucnt = {};

        ['login', 'bill'].forEach(t => {
            var m = {};
            result[t].forEach(row => {
                row.uids.forEach(uid => {
                    m[uid] = true;
                });
            });
            sum_ucnt[t] = Object.keys(m).length;
        });

        // merge
        {
            var m = {}
            result.login.forEach(row => {
                m[row._id] = row.uids.length;
            });

            result.bill.forEach(row => {
                var n = m[row._id];
                row.login_ucnt = n ? n : 0;
            });
        }

        // tab
        var tab = {};

        // header
        tab.header = ["日期", "活跃人数", "充值人数", "充值金额", "arppu", "活跃arpu", "付费率"];

        // body
        tab.body = result.bill.map(row => [
            row._id,
            row.login_ucnt,
            row.uids.length,
            row.amt,
            (row.amt / row.uids.length).toFixed(2),
            (row.amt / row.login_ucnt).toFixed(2),
            (row.uids.length * 100 / row.login_ucnt).toFixed(2) + '%',
        ]);

        var sum_amt = tab.body.map(v => v[3]).reduce((a, v) => a + v, 0);

        tab.body.forEach(v => {
            v[3] = v[3] + ` (${gtab.toCNY(Number(v[3])).toFixed(2)} CNY)`;
            v[4] = v[4] + ` (${gtab.toCNY(Number(v[4])).toFixed(2)} CNY)`;
            v[5] = v[5] + ` (${gtab.toCNY(Number(v[5])).toFixed(2)} CNY)`;
        });

        // sum
        tab.body.push([
            "合计",
            sum_ucnt.login,
            sum_ucnt.bill,
            sum_amt + ` (${Math.floor(gtab.toCNY(sum_amt))} CNY)`,
            (sum_amt / sum_ucnt.bill).toFixed(2) + ` (${gtab.toCNY(sum_amt / sum_ucnt.bill).toFixed(2)} CNY)`,
            (sum_amt / sum_ucnt.login).toFixed(2) + ` (${gtab.toCNY(sum_amt / sum_ucnt.login).toFixed(2)} CNY)`,
            (sum_ucnt.bill * 100 / sum_ucnt.login).toFixed(2) + '%',
        ]);

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
