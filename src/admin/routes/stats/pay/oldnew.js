
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
        tab.name = '新老充值';
        stats.render(req, res, tab);
    }

    // ----------------------

    var tz = new Date().getTimezoneOffset() * 60 * 1000;

    // ---------------------

    var pipeline = [
        {$match: {
            bts: {$gte: q.d_start, $lt: q.d_end},
            amt: {$gt: 0},
        }},
        {$project: {
            bts: {$subtract: ["$bts", tz]},
            day: 1,
            amt: 1,
        }},
        {$group: {
            _id:     {$dateToString: {format: "%Y-%m-%d", date: "$bts"}},
            n:       {$sum: 1},
            n_old:   {$sum: {$cond: [{$ne: ["$day", 1]}, 1, 0]}},
            n_new:   {$sum: {$cond: [{$eq: ["$day", 1]}, 1, 0]}},
            amt:     {$sum: "$amt"},
            amt_old: {$sum: {$cond: [{$ne: ["$day", 1]}, "$amt", 0]}},
            amt_new: {$sum: {$cond: [{$eq: ["$day", 1]}, "$amt", 0]}},
        }},
        {$sort: {_id: -1}},
    ];

    stats.filter(pipeline, q);

    try {
        let db = dbpool.get('stats');

        let docs = await db.collection('bill').aggregate(pipeline).toArray();

        // tab
        var tab = {};

        // header
        tab.header = ["日期", "老用户充值人数", "老用户充值金额", "新用户充值人数", "新用户充值金额", "充值人数", "充值金额"];

        // body
        tab.body = docs.map(row => [
                row._id,
                row.n_old,
                row.amt_old + ` (${gtab.toCNY(row.amt_old).toFixed(2)} CNY)`,
                row.n_new,
                row.amt_new + ` (${gtab.toCNY(row.amt_new).toFixed(2)} CNY)`,
                row.n,
                row.amt + ` (${gtab.toCNY(row.amt).toFixed(2)} CNY)`,
        ]);

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
