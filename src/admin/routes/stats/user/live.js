
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
        tab.name = '用户留存';
        stats.render(req, res, tab);
    }

    // ----------------------

    var tz = new Date().getTimezoneOffset() * 60 * 1000;

    // ---------------------

    var days = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 30, 40, 50, 60, 70, 80, 90];

    // ---------------------

    var obj1 = {
        _id: {$dateToString: {format: "%Y-%m-%d", date: "$cts"}},
    };
    days.forEach(v => {
        obj1[`d${v}`] = {$sum: {$cond: [{$eq: ["$day", v]}, 1, 0]}};

    });

    // ---------------------

    var pipeline = [
        {$match: {
            cts: {$gte: q.d_start, $lt: q.d_end},
            day: {$lte: 90},
        }},
        {$project: {
            cts: {$subtract: ["$cts", tz]},
            day: 1,
        }},
        {$group: obj1},
        {$sort: {_id: -1}},
    ];

    stats.filter(pipeline, q);

    // ---------------------

    try {
        let db = dbpool.get('stats');

        let docs = await db.collection('login').aggregate(pipeline).toArray();

        // tab
        var tab = {};

        // header
        tab.header = ['日期', '新增'];
        days.shift();
        days.forEach(v => {
            tab.header.push(`${v}日留存`);
        });

        // body
        tab.body = docs.map(row => {
            var r = [row._id, row.d1];
            days.forEach(v => {
                r.push((row[`d${v}`] * 100 / row.d1).toFixed(2) + '%');
            });
            return r;
        });

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
