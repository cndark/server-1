
var express = require('express');
var router  = express.Router();

var dbpool  = require('../../../lib/dbpool');
var stats   = require('../stats');

// ============================================================================

router.get('/', _A_(async (req, res) => {
    var q = stats.make_q(req.query);

    // render func
    var render = (tab) => {
        tab = tab || {};
        tab.name = '用户概览';
        stats.render(req, res, tab);
    }

    // ----------------------

    var tz = new Date().getTimezoneOffset() * 60 * 1000;

    // ---------------------

    try {
        let db = dbpool.get('stats');

        let fs = {
            create: async () => {
                let pipeline = [
                    {$match: {
                        cts: {$gte: q.d_start, $lt: q.d_end},
                    }},
                    {$project: {
                        cts: {$subtract: ["$cts", tz]},
                        day: 1,
                    }},
                    {$group: {
                        _id: {$dateToString: {format: "%Y-%m-%d", date: "$cts"}},
                        n:   {$sum: {$cond: [{$eq: ["$day", 1]}, 1, 0]}},
                    }},
                ];

                stats.filter(pipeline, q);

                return await db.collection('login').aggregate(pipeline).toArray();
            },

            active: async () => {
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
                ];

                stats.filter(pipeline, q);

                return await db.collection('login').aggregate(pipeline).toArray();
            },

            online: async () => {
                let pipeline = [
                    {$match: {
                        ts: {$gte: q.d_start, $lt: q.d_end},
                    }},
                    {$project: {
                        ts: {$subtract: ["$ts", tz]},
                        avg:  1,
                        peek: 1,
                    }},
                    {$group: {
                        _id:  {$dateToString: {format: "%Y-%m-%d", date: "$ts"}},
                        avg:  {$avg: "$avg"},
                        peek: {$max: "$peek"},
                    }},
                ];

                stats.filter(pipeline, q);
                delete pipeline[0].$match.sdk; // ignore sdk

                return await db.collection('online').aggregate(pipeline).toArray();
            },
        };

        let result = await Promise.all(['create', 'active', 'online'].map(k => fs[k]()));
        result.create = result[0];
        result.active = result[1];
        result.online = result[2];

        // calc active-count for sum
        var sum_active = {};

        result.active.forEach(row => {
            row.uids.forEach(uid => {
                sum_active[uid] = true;
            });
        });
        sum_active = Object.keys(sum_active).length;

        // tab
        var tab = {};

        // header
        tab.header = ["日期", "新增", "活跃", "平均在线", "最高在线"];

        // body
        tab.body = {};

        var merge = function (i, rs, f) {
            for (var row of rs) {
                if (!tab.body[row._id])
                    tab.body[row._id] = [row._id, 0, 0];

                tab.body[row._id][i] = f(row);
            }
        };

        merge(1, result.create, row => row.n);
        merge(2, result.active, row => row.uids.length);
        merge(3, result.online, row => Number(row.avg.toFixed(0)));
        merge(4, result.online, row => row.peek);

        tab.body = Object.keys(tab.body).sort((a, b) => b.localeCompare(a)).map(v => tab.body[v]);

        // sum
        tab.body.push([
            "合计",
            tab.body.map(v=>v[1]).reduce((a, v)=>a + v, 0),
            sum_active,
            (tab.body.map(v=>v[3]).reduce((a, v)=>a + v, 0) / (tab.body.length == 0 ? 1 : tab.body.length)).toFixed(0),
            tab.body.map(v=>v[4]).reduce((a, v)=>a > v ? a : v, 0),
        ]);

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
