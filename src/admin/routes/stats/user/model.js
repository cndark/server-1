
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
        tab.name = '设备型号';
        stats.render(req, res, tab);
    }

    // ----------------------

    var pipeline = [
        {$match: {
            ts: {$gte: q.d_start, $lt: q.d_end},
        }},
        {$group: {
            _id: "$model",
            n:   {$sum: 1},
        }},
        {$sort: {n: -1}},
    ];

    stats.filter(pipeline, q);
    delete pipeline[0].$match.area; // ignore area
    delete pipeline[0].$match.svr;  // ignore svr

    // ---------------------

    try {
        let db = dbpool.get('stats');

        let docs = await db.collection('install').aggregate(pipeline).toArray();

        // tab
        var tab = {};

        // header
        tab.header = ["型号", "数量"];

        // body
        tab.body = docs.map(row => [row._id, row.n]);

        // sum
        tab.body.push([
            "合计",
            docs.map(v => v.n).reduce((a, v) => a + v, 0),
        ]);

        // render
        render(tab);

    } catch {
        render();
    }
}));

// ============================================================================

module.exports = router;
