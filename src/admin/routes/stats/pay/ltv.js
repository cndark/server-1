
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
        tab.name = '付费LTV';
        stats.render(req, res, tab);
    }

    // ----------------------

    var tz = new Date().getTimezoneOffset() * 60 * 1000;

    // ---------------------

    var days = [];
    for (let i = 1; i <= 30; i++)
        days.push(i);

    // ---------------------

    var obj1 = {
        _id:  {$dateToString: {format: "%Y-%m-%d", date: "$cts"}},
        ucnt: {$sum: {$cond: [{$eq: ["$day", 1]}, 1, 0]}},
    };
    days.forEach(v => {
        obj1[`d${v}`] = {$sum: {$cond: [{$eq: ["$day", v]}, "$amt", 0]}};
    });

    // ---------------------

    var pipeline = [
        {$match: {
            cts: {$gte: q.d_start, $lt: q.d_end},
            day: {$lte: 30},
        }},
        {$project: {
            cts: {$subtract: ["$cts", tz]},
            day: 1,
            amt: 1,
        }},
        {$group: obj1},
        {$sort: {_id: -1}},
    ];

    // ---------------------

    stats.filter(pipeline, q);

    // ---------------------

    try {
        let db = dbpool.get('stats');

        let docs = await db.collection('bill').aggregate(pipeline).toArray();

        // tab
        var tab = {};

        // header
        tab.header = ['日期'];
        days.forEach(v => {
            tab.header.push(`ltv${v}`);
        });

        // body
        var footer = [];
        for (let i = 1; i <= 30; i++)
            footer[i] = {sum: 0, ucnt: 0};

        var footer_max_day = 0;

        tab.body = docs.map(row => {
            var r       = [row._id];
            var sum     = 0;
            var max_day = Math.floor((Date.now() - Date.fromString(row._id).getTime()) / (86400 * 1000)) + 1;

            days.forEach(v => {
                sum += row[`d${v}`];

                if (v > max_day) {
                    r.push('-');
                } else {
                    r.push((sum / row.ucnt).toFixed(2) + ` (${gtab.toCNY(sum / row.ucnt).toFixed(2)}_CNY)`);
                    if (v > footer_max_day) footer_max_day = v;
                }

                footer[v].sum  += sum;
                footer[v].ucnt += row.ucnt;
            });

            return r;
        });

        // footer
        var r = ['合计'];
        days.forEach(v => {
            if (v > footer_max_day)
                r.push('-');
            else
                r.push((footer[v].sum / footer[v].ucnt).toFixed(2) + ` (${gtab.toCNY(footer[v].sum / footer[v].ucnt).toFixed(2)}_CNY)`);
        });

        tab.body.push(r);

        // css
        tab.css = `
            table td:nth-child(1) {
                min-width: 100px;
            }
        `;

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
